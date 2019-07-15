pragma solidity >=0.4.22;

/// @title Proof of Authority Whitelist Proof of Concept
/// @author Jon Knight
/// @author Mosaic Networks
/// @notice Copyright Mosaic Networks 2019, released under the MIT license

contract POA_Genesis {
/// @notice Event emitted when the vote was reached a decision
/// @param _nominee The address of the nominee
/// @param _yesVotes The total number of yes votes cast for the nominee to date
/// @param _noVotes The total number of no votes cast for the nominee to date
/// @param _accepted The decision, true for added to the whitelist, false for rejected
    event NomineeDecision(
        address indexed _nominee,
        uint  _yesVotes,
        uint _noVotes,
        bool indexed _accepted
    );

/// @notice Event emitted when a nominee vote is cast
/// @param _nominee The address of the nominee
/// @param _voter The address of the person who cast the vote
/// @param _yesVotes The total number of yes votes cast for the nominee to date
/// @param _noVotes The total number of no votes cast for the nominee to date
/// @param _accepted The vote, true for accept, false for rejected
    event NomineeVoteCast(
        address indexed _nominee,
        address indexed _voter,
        uint  _yesVotes,
        uint _noVotes,
        bool indexed _accepted
    );





/// @notice Event emitted when a nominee is proposed
/// @param _nominee The address of the nominee
/// @param _proposer The address of the person who proposed the nominee
    event NomineeProposed(
        address indexed _nominee,
        address indexed _proposer
    );



/// @notice Event emitted to announce a moniker
/// @param _address The address of the user
/// @param _moniker The moniker of the user
    event MonikerAnnounce(
        address indexed _address,
        bytes32 indexed _moniker
    );



   struct WhitelistPerson {
      address person;
      uint  flags;
   }

   struct NomineeVote {
      address voter;
      bool  accept;
   }

   struct NomineeElection{
      address nominee;
      address proposer;
      uint yesVotes;
      uint noVotes;
      mapping (address => NomineeVote) vote;
      address[] yesArray;
      address[] noArray;
   }

   mapping (address => WhitelistPerson) public whiteList;
   uint whiteListCount;
   address[] whiteListArray;
   mapping (address => NomineeElection) nomineeList;
   address[] nomineeArray;
   mapping (address => bytes32) monikerList;


//GENERATED GENESIS BEGIN

    // This code should be removed as part of the monetcli network compile  /monetcli wizard process. 
    // The code within here is included so that the unamended contract is viable. 

    address constant initWhitelist1 = 0x1234567890123456789012345678901234567890;
    address constant initWhitelist2 = 0x2345678901234567890123456789012345678901;
    address constant initWhitelist0 = 0x3456789012345678901234567890123456789012;

    bytes32 constant initWhitelistMoniker1 = "Tom";
    bytes32 constant initWhitelistMoniker2 = "Dick";
    bytes32 constant initWhitelistMoniker0 = "Harry";

   /// @notice function to apply the genesis white list to the whitelist. This will be generated.
    function processGenesisWhitelist() private
    {
        addToWhitelist(initWhitelist1, initWhitelistMoniker1);
        addToWhitelist(initWhitelist2, initWhitelistMoniker2);
        addToWhitelist(initWhitelist0, initWhitelistMoniker0);
    }


   /// @notice function to check if an address is on the intial genesis block white list
   /// @param _address the address to be checked
   /// @return a boolean value, indicating if _address is on the white list
    function isGenesisWhitelisted(address _address) pure private returns (bool)
    {
        return (  ( initWhitelist1 == _address ) ||  ( initWhitelist2 == _address ));
    }

    // Code down to here should be replaced. 

//GENERATED GENESIS END

// There is no constructor for the genesis block

   /// @notice Constructor builds the white list with just the contract owner in it
   /// @param _moniker the name of the contract owner as shown to other users in the wallet
    constructor(bytes32 _moniker) public {
        addToWhitelist(address(uint160(msg.sender)), _moniker);
        processGenesisWhitelist();
    }


   /// @notice This is a constructor replacement for contracts placed directly in the genesis block. This is necessary because the constructor does not run in that instance.
    function init () public payable checkAuthorisedModifier(msg.sender)
    {
    	processGenesisWhitelist();
    }


   /// @notice Modifier to check if a sender is on the white list.
    modifier checkAuthorisedModifier(address _address)
    {
        if (whiteListCount == 0){
            require(isGenesisWhitelisted(_address), "Not authorised");
            // This is a modifier on a payable transaction so we can initialise everything.
            processGenesisWhitelist();
        }
        require(isWhitelisted(_address), "Not authorised");
        _;
    }


   /// @notice Function exposed for Babble Join authority
   function checkAuthorised(address _address) public view returns (bool)
   {  // needs check on whitelist to allow original validators to be booted.

       return ((isWhitelisted(_address)) || ((whiteListCount == 0)&&(isGenesisWhitelisted(_address)))  );
   }

   /// @notice Function exposed for Babble Join authority wraps checkAuthorised
   function checkAuthorisedPublicKey(bytes32  _publicKey) public view returns (bool)
   {
      return checkAuthorised(address(uint160(uint256(keccak256(abi.encodePacked(_publicKey))))));

//    This version works in Solidity 0.4.x, but the extra intermediate steps are required in 0.5.x
//      return checkAuthorised(address(keccak256(abi.encodePacked(_publicKey))));
   }

   /// @notice wrapper function to check if an address is on the nominee list
   /// @param _address the address to be checked
   /// @return a boolean value, indicating if _address is on the nominee list
    function isNominee(address _address) private view returns (bool)
    {
        return (nomineeList[_address].nominee != address(0));
    }

   /// @notice wrapper function to check if an address is on the white list
   /// @param _address the address to be checked
   /// @return a boolean value, indicating if _address is on the white list
    function isWhitelisted(address _address) private view returns (bool)
    {
        return (whiteList[_address].person != address(0));
    }


    /// @notice private function to add user directly to the whitelist. Used to process the Genesis Whitelist.
    function addToWhitelist(address _address, bytes32 _moniker) private {

        if (! isWhitelisted(_address))   // prevent duplicate whitelist entries
        {
           whiteList[_address] = WhitelistPerson(_address, 0);
           whiteListCount++;
           whiteListArray.push(_address);
           monikerList[_address] = _moniker;
           emit MonikerAnnounce(_address,_moniker);
           emit NomineeDecision(_address, 0, 0, true);  // zero vote counts because there was no vote
        }
    }


   /// @notice Add a new entry to the nominee list
   /// @param _nomineeAddress the address of the nominee
   /// @param _moniker the moniker of the new nominee as displayed during the voting process
    function submitNominee (address _nomineeAddress, bytes32 _moniker) public payable checkAuthorisedModifier(msg.sender)
    {
        nomineeList[_nomineeAddress] = NomineeElection({nominee: _nomineeAddress, proposer: msg.sender,
                    yesVotes: 0, noVotes: 0, yesArray: new address[](0),noArray: new address[](0) });
        nomineeArray.push(_nomineeAddress);
        monikerList[_nomineeAddress] = _moniker;
        emit NomineeProposed(_nomineeAddress,  msg.sender);
        emit MonikerAnnounce(_nomineeAddress, _moniker);
    }


    ///@notice Cast a vote for a nominator. Can only be run by people on the whitelist.
    ///@param _nomineeAddress The address of the nominee
    ///@param _accepted Whether the vote is to accept (true) or reject (false) them.
    ///@return returns true if the vote has reached a decision, false if not
    ///@return only meaningful if the other return value is true, returns true if the nominee is now on the whitelist. false otherwise.
    function castNomineeVote(address _nomineeAddress, bool _accepted) public payable checkAuthorisedModifier(msg.sender) returns (bool decided, bool voteresult){

        decided = false;
        voteresult = false;

//      Check if open nominee, other checks redundant
        if (isNominee(_nomineeAddress)) {


//      Check that this sender has not voted before. Initial config is no redos - so just reject
            if (nomineeList[_nomineeAddress].vote[msg.sender].voter == address(0)) {
                // Vote is valid. So lets cast the Vote
                nomineeList[_nomineeAddress].vote[msg.sender] = NomineeVote({voter: msg.sender, accept: _accepted });

                // Amend Totals
                if (_accepted)
                {
                    nomineeList[_nomineeAddress].yesVotes++;
                    nomineeList[_nomineeAddress].yesArray.push(msg.sender);
                } else {
                    nomineeList[_nomineeAddress].noVotes++;
                    nomineeList[_nomineeAddress].noArray.push(msg.sender);
                }

                emit NomineeVoteCast(_nomineeAddress, msg.sender,nomineeList[_nomineeAddress].yesVotes,
                      nomineeList[_nomineeAddress].noVotes, _accepted);

                // Check to see if enough votes have been cast for a decision
                (decided, voteresult) = checkForNomineeVoteDecision(_nomineeAddress);
            }
        }
        else
        {   // Not a nominee, so set decided to true
            decided = true;
        }

        // If decided, check if on whitelist
        if (decided) {
            voteresult = isWhitelisted(_nomineeAddress);
        }
        return (decided, voteresult);
    }

    ///@notice This function encapsulates the logic for determining if there are enough votes for a definitive decision
    ///@param _nomineeAddress The address of the NomineeElection
    ///@return returns true if the vote has reached a decision, false if not
    ///@return only meaningful if the other return value is true, returns true if the nominee is now on the whitelist. false otherwise.

    function checkForNomineeVoteDecision(address _nomineeAddress) private returns (bool decided, bool voteresult)
    {
        NomineeElection memory election = nomineeList[_nomineeAddress];
        decided = false;
        voteresult = false;


        if (election.noVotes > 0)  // Someone Voted No
        {
            declineNominee(election.nominee);
            decided = true;
            voteresult = false;
        }
        else
        {
            // Requires unanimous approval
            if(election.yesVotes >= whiteListCount)
            {
                acceptNominee(election.nominee);
                decided = true;
                voteresult = true;
            }
        }

        if (decided)
        {
            emit NomineeDecision(_nomineeAddress, election.yesVotes, election.noVotes, voteresult);
        }
        return (decided, voteresult);
    }


    ///@notice This private function adds the accepted nominee to the whitelist.
    ///@param _nomineeAddress The address of the nominee being added to the whitelist
    function acceptNominee(address _nomineeAddress) private
    {
        if (! isWhitelisted(_nomineeAddress))  // avoid re-adding and corrupting the whiteListCount
        {
          whiteList[_nomineeAddress] = WhitelistPerson(_nomineeAddress, 0);
          whiteListArray.push(_nomineeAddress);
          whiteListCount++;
        }
    // Remove from nominee list
       removeNominee(_nomineeAddress);
    }


    ///@notice This private function adds the removes a user from the whitelist. Not currently used.
    ///@param _address The address of the nominee being removed to the whitelist
    function deWhiteList(address _address) private
    {
        if(isWhitelisted(_address))
        {
            delete(whiteList[_address]);
            whiteListCount--;

            for (uint i = 0; i<whiteListArray.length; i++) {
                if (whiteListArray[i] == _address)
                {  // Replace item to be removed with the last item. Then remove last item.
                    whiteListArray[i] == whiteListArray[whiteListArray.length - 1];
                    delete whiteListArray[whiteListArray.length - 1];
                break;
                }
            }
        }
    }


// Deline nominee from the nomineeList

    ///@notice This private function removes the declined nominee from the nominee list.
    ///@param _nomineeAddress The address of the nominee being removed from the nominee list
    function declineNominee(address _nomineeAddress) private
    {
         removeNominee(_nomineeAddress);
    }


    ///@notice This private function removes the declined nominee from the nominee list.
    ///@param _nomineeAddress The address of the nominee being removed from the nominee list
    function removeNominee(address _nomineeAddress) private
    {
// Remove from Mapping
        delete(nomineeList[_nomineeAddress]);

			for (uint i = 0; i<nomineeArray.length; i++) {
				if (nomineeArray[i] == _nomineeAddress)
				{  // Replace item to be removed with the last item. Then remove last item.
					nomineeArray[i] == nomineeArray[nomineeArray.length - 1];
					delete nomineeArray[nomineeArray.length - 1];
                  break;
				}
			}
    }



// Information Section.

	function getNomineeElection(address _address) public view returns (address nominee, address proposer, uint yesVotes, uint noVotes)
	{
		return (nomineeList[_address].nominee, nomineeList[_address].proposer, nomineeList[_address].yesVotes, nomineeList[_address].noVotes);
	}


// Array Section. Functions to support Arrays.



	function getNomineeCount() public view returns (uint count)
	{
		return (nomineeArray.length);
	}


   function getNomineeAddressFromIdx(uint idx) public view returns (address NomineeAddress)
	{
		require (idx < nomineeArray.length, "Requested address is out of range.");
		return (nomineeArray[idx]);
	}


	function getNomineeElectionFromIdx(uint idx) public view returns (address nominee, address proposer, uint yesVotes, uint noVotes)
	{
		return (getNomineeElection(getNomineeAddressFromIdx(idx))) ;

	}


	function getWhiteListCount() public view returns (uint count)
	{
		return (whiteListArray.length);
	}


   function getWhiteListAddressFromIdx(uint idx) public view returns (address WhiteListAddress)
	{
		require (idx < whiteListArray.length, "Requested address is out of range.");
		return (whiteListArray[idx]);
	}


	function getYesVoteCount(address _nomineeAddress)  public view returns (uint count)
	{
		return (nomineeList[_nomineeAddress].yesArray.length);
	}

	function getYesVoterFromIdx(address _nomineeAddress, uint _idx)  public view returns (address voter)
	{
		require (_idx < nomineeList[_nomineeAddress].yesArray.length, "Requested address is out of range.");
		return (nomineeList[_nomineeAddress].yesArray[_idx]);
	}


	function getNoVoteCount(address _nomineeAddress)  public view returns (uint count)
	{
		return (nomineeList[_nomineeAddress].noArray.length);
	}

	function getNoVoterFromIdx(address _nomineeAddress, uint _idx) public view returns (address voter)
	{
		require (_idx < nomineeList[_nomineeAddress].noArray.length, "Requested address is out of range.");
		return (nomineeList[_nomineeAddress].noArray[_idx]);
	}

	function getMoniker(address _address) public view returns (bytes32 moniker)
	{
		return (monikerList[_address]);
	}

        function getCurrentNomineeVotes(address _address) public view returns (uint yes, uint no)
    {
       if (! isNominee(_address)) {return (yes, no);}
        return (nomineeList[_address].yesVotes,nomineeList[_address].noVotes);
    }

}
