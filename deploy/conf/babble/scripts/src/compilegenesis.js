const minimist = require('minimist');
const program = require('commander');
const fs = require('fs');
const path = require('path');
const createKeccakHash = require('keccak')
const shelljs = require('shelljs');


program
	.version('0.0.2', "-v, --version")
   .option('-p, --precompfile <file>', 'pregenesis.json file')
   .option('-o, --output-dir <directory>', 'output directory')
   .parse(process.argv);



if ( program.precompfile == undefined) { program.precompfile = 'pregenesis.json';}
if ( program.outputDir == undefined) { program.outputDir = '/tmp';}



// Load the pregenesis.json file
function loadFile(_pcFile)
{
    try
    {
    	rawprecomp = fs.readFileSync(_pcFile);
    } catch (err)
    {
      if (err.code === 'ENOENT') {
        console.log('File not found!');
        console.log( _pcFile );
        return undefined; //TODO maybe make this more sensible
      } else {
        throw err;
      }
    }
    
    return rawprecomp;    
}



function loadFileUTF(_pcFile)
{
    try
    {
    	rawprecomp = fs.readFileSync(_pcFile, "utf8");
    } catch (err)
    {
      if (err.code === 'ENOENT') {
        console.log('File not found!');
        console.log(_pcFile);
        return undefined; //TODO maybe make this more sensible
      } else {
        throw err;
      }
    }
    
    return rawprecomp;    
}


function writeFile(_file, _data)
{

   let fullPath = path.join(program.outputDir, _file);

   console.log("Writing file", fullPath);
   
	fs.writeFileSync(fullPath, _data);

   console.log("Written file", fullPath);
}



function tidyAddress(address)
{
  address = address.toLowerCase().replace('0x', '');
  return address;
}

function tidyAddressUpper(address)
{
  address = '0X' + address.toUpperCase().replace('0X', '');
  return address;
}


function tidyAddressSol(address)
{

  address = address.toLowerCase().replace('0x', '')
  var hash = createKeccakHash('keccak256').update(address).digest('hex')
  var ret = '0x'

  for (var i = 0; i < address.length; i++) {
    if (parseInt(hash[i], 16) >= 8) {
      ret += address[i].toUpperCase()
    } else {
      ret += address[i]
    }
  }

  return ret
}


function generateSolidityHardCodings(_preAuthorised)
{
  let rtn = '';

  if ( ( _preAuthorised ) && ( Array.isArray(_preAuthorised))  )  
  {
		var consts = [];
		var addTo = [];
		var checks = [];

		for (var i = 0; i < _preAuthorised.length; i++)
		{
			let person = _preAuthorised[i];

			if ((person.address) && (person.moniker))
			{
			    let tidiedaddress = tidyAddressSol(person.address); 
				 consts.push('    address constant initWhitelist'+i+' = '+tidiedaddress +';'); 	
				 consts.push('    bytes32 constant initWhitelistMoniker'+i+' = "'+person.moniker +'";'); 	

             addTo.push('        addToWhitelist(initWhitelist'+i+', initWhitelistMoniker'+i+');');
    
             checks.push(' ( initWhitelist'+i+' == _address ) ');

			}
		}
  }

//  console.log(checks.join('||'));

  rtn = " //GENERATED GENESIS BEGIN \n " +
		  " \n" +	
		  consts.join("\n")+
		  " \n" +	
		  " \n" +	
		  " \n" +	
        "    function processGenesisWhitelist() private \n" +
        "    { \n" +
        addTo.join("\n")+
		  " \n" +	
        "    } \n" +
		  " \n" +	
		  " \n" +	
        "    function isGenesisWhitelisted(address _address) pure private returns (bool) \n"+
        "    { \n"+
        "        return ( "+checks.join('||') + "); \n"+
        "    } \n"+

		  " \n" +	
        " //GENERATED GENESIS END \n " ;


   return rtn;
}




function processContract(_output, _contract, _populate_alloc, _populate_poa)
{
	
 //  console.dir(_contract, { depth: 6, colors: true });

   if (_contract.address)
	{
		if (! _output.alloc) { _output.alloc = {} ; } 
		
		let tidiedaddress = tidyAddressUpper(_contract.address); 
		
		if (! _output.alloc[tidiedaddress] )   // make sure the address node exists in the alloc hierarchy
      {
			_output.alloc[tidiedaddress] = {};
		}  	

		if (_contract.balance) { _output.alloc[tidiedaddress].balance = _contract.balance;}

		if (_contract.filename)
		{

                         
			 let rawsol = loadFileUTF(__dirname + "/../../conf/poa/"+ _contract.filename);
			 if (rawsol) 
			 {
             let newcode = generateSolidityHardCodings(_contract.preauthorised);

//				console.log('matches', rawsol.match(/\/\/GENERATED GENESIS BEGIN[\S\s]*GENERATED GENESIS END/g));

				 let newsol = rawsol.replace(/\/\/GENERATED GENESIS BEGIN[\S\s]*GENERATED GENESIS END/, newcode);	

				 return newsol;
//				 console.log(newsol);

			 }
			 else
			 {
				 console.log('Solidity File not found');	
			 } 	
		}

		
//		console.log(generateSolidityHardCodings(_contract.preauthorised));



	}
	else
	{
		console.log('This contract does not have an address specified and thus has not been applied to the genesis block');
   }

	return undefined;
}



function processPreGenesisFile(_pcFile, _populate_alloc, _populate_poa)
{
//     console.log(_pcFile);
    let rawprecomp = loadFile(_pcFile)
    if (! rawprecomp) { return ;}

    let precomp = JSON.parse(rawprecomp);
    let output = JSON.parse(rawprecomp); // least bad way to clone the object.
    



    
    if ( precomp.precompiler)
    {

//      Commented out to see if it still errors
//    	delete(output.precompiler);
    
    	if ( (precomp.precompiler.contracts) && ( Array.isArray(precomp.precompiler.contracts) ))
    	{
       let contracts = precomp.precompiler.contracts;
			 for (var i=0; i<contracts.length ; i++)
			 {
					let contract = processContract(output, contracts[i], _populate_alloc, _populate_poa);
               let contractfilename = 'contract'+i+'.sol';
             
               if ( _populate_alloc && (contracts[i].authorising) )  // poa section is authorising by definition
               {		
                   output.alloc[tidyAddressUpper(contracts[i].address)].authorising = true ;
               }

               writeFile(contractfilename, contract);					

//					console.log(__dirname+'/../compile_solidity.sh '+ program.outputDir + ' ' + contractfilename );
					shelljs.exec(__dirname+'/../compile_solidity.sh --output-dir '+ program.outputDir + ' --contract ' + path.join(program.outputDir, contractfilename ));

					if (contracts[i].contractname)
					{
            let bytecode = loadFileUTF(path.join(program.outputDir, contracts[i].contractname+ '.bin-runtime'));
            let abicode = loadFileUTF(path.join(program.outputDir, contracts[i].contractname+ '.abi'));
            let abijson = JSON.parse(abicode);

						if (bytecode)    //TODO need to check address is set, but I'll miss my train`
						{
                if (_populate_alloc)
                {  
                    output.alloc[tidyAddressUpper(contracts[i].address)].code = bytecode ; //  "0X"+upper(bytecode);
                }
                if (_populate_poa)   // NB is it a feature that is multiple contracts are defined, only the latest is processed.
                {
                    output.poa = {};
                    output.poa.address =  tidyAddressUpper(contracts[i].address);
                    output.poa.abi = abicode;
                    output.poa.code = bytecode;
                }
						}
						else
						{

							console.log ('Bytecode not found',path.join(program.outputDir, contracts[i].contractname+ '.bin-runtime')); 
						}
							
					}
					

          }  
    
    	}
    	else
    	{   // No contracts defined
    		console.log('Precompiler section present, but no contracts defined.');
       }	
    
    }
    else
    {  // Pure genesis.json file with no contract / shenanigans
    	console.log('Nothing to precompile'); 	
    }

// If not populating the alloc section, we do not need the precompiler section
   if (! _populate_alloc ) {delete(output.precompiler);}

 //     console.log('');

 //    console.dir(output, { depth: 6, colors: true });
	  writeFile('genesis.json' , JSON.stringify(output, null, 3));	    
}



// There is a migration from placing the contract in the alloc section to placing it 
// in the poa section. As the code is quite involved, we set parameters here to control the 
// output. 

let populate_alloc = false;   // Populate alloc section
let populate_poa = true;     // Populate POA section


processPreGenesisFile(program.precompfile, populate_alloc, populate_poa);

