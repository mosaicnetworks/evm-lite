pragma solidity >=0.4.25;

/// @title Proof of Authority Whitelist Proof of Concept
/// @author Jon Knight
/// @author Mosaic Networks
/// @notice This is a proof of concept and is not production ready


contract MonetPOA {


/// @notice Event emitted when user added to whitelist without taking a vote
    event NodeAuthorised(
        address indexed _authorisee,
        address indexed _authorisor
    );



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
        uint8 _progress,
        bool indexed _accepted
    );



   uint constant ERR_ALREADY_VOTED = 1;
   uint constant ERR_NOT_ON_NOMINEELIST = 2;


    event DebugMessageEvent(
      address indexed _address,
      uint errcode1,
      address erraddress
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



   /// @notice Constructor builds the white list with just the contract owner in it
   /// @param _moniker the name of the contract owner as shown to other users in the wallet
   //  constructor(string _moniker) public ;
   /// @notice Modifier to check if a sender is on the white list.
    modifier checkAuthorised(address _address)
    {
        _;
    }


     
   /// @notice Add a new entry to the nominee list
   /// @param _nomineeAddress the address of the nominee
   /// @param _moniker the moniker of the new nominee as displayed during the voting process    
    function submitNominee (address _nomineeAddress, bytes32 _moniker) payable public;



    ///@notice Cast a vote for a nominator. Can only be run by people on the whitelist. 
    ///@param _nomineeAddress The address of the nominee
    ///@param _accepted Whether the vote is to accept (true) or reject (false) them. 
    ///@return returns true if the vote has reached a decision, false if not
    ///@return only meaningful if the other return value is true, returns true if the nominee is now on the whitelist. false otherwise.
    function castNomineeVote(address _nomineeAddress, bool _accepted) public payable returns (bool, bool);

}    




contract POA_Event2 is MonetPOA {

   

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
   }


   mapping (address => WhitelistPerson) public whiteList; 
   uint whiteListCount;    
    



// Cannot iterate through a mapping, so keep an array in place too. 
   mapping (address => NomineeElection) nomineeList;   



   /// @notice Constructor builds the white list with just the contract owner in it
   /// @param _moniker the name of the contract owner as shown to other users in the wallet
    constructor(bytes32 _moniker) public {
       whiteList[address(uint160(msg.sender))] = WhitelistPerson(msg.sender, 0);
       whiteListCount = 1;
       emit MonikerAnnounce(msg.sender,_moniker);
       emit NodeAuthorised(msg.sender, msg.sender);
    }


   /// @notice Modifier to check if a sender is on the white list.
    modifier checkAuthorised(address _address)
    {
        require(isWhitelisted(_address));

        _;
    }


   /// @notice wrapper function to check if an address is on the nominee list
   /// @param _address the address to be checked
   /// @return a boolean value, indicating if _address is on the nominee list
    function isNominee(address _address) public view returns (bool)
    {
        return (nomineeList[_address].nominee != address(0));
    }


   /// @notice wrapper function to check if an address is on the white list
   /// @param _address the address to be checked
   /// @return a boolean value, indicating if _address is on the white list
    function isWhitelisted(address _address) public view returns (bool)
    {
        return (whiteList[_address].person != address(0));
    }



//TODO - Testing Functions, not for production use BEGIN

   /// @notice Testing function to return the length of the whitelist.   
    function getWhiteListCount() public view returns (uint)
   {
      return whiteListCount;
   }



   /// @notice Get Current Status of a Vote
    function getVoteCounts(address _nomineeAddress) public view returns (uint, uint, uint, uint)
   {
      if ( ! isNominee(_nomineeAddress) ) 
      {
         return(0, whiteListCount, 0, 0);
      }
      else
      {
         return(1, whiteListCount,nomineeList[_nomineeAddress].yesVotes, nomineeList[_nomineeAddress].noVotes); 
      }

   }


   /// @notice Cheat and add another node to the whitelist

    function authoriseNodeWithoutVote(address _address, bytes32 _moniker) payable public {
       whiteList[_address] = WhitelistPerson(_address, 0);
       whiteListCount++;
       emit MonikerAnnounce(_address, _moniker);
       emit NodeAuthorised(_address, msg.sender);
    }

     
//TODO - Testing Functions, not for production use END



   /// @notice Add a new entry to the nominee list
   /// @param _nomineeAddress the address of the nominee
   /// @param _moniker the moniker of the new nominee as displayed during the voting process    
    function submitNominee (address _nomineeAddress, bytes32 _moniker) payable public checkAuthorised(msg.sender)
    {    
        nomineeList[_nomineeAddress] = NomineeElection({nominee: _nomineeAddress, proposer: msg.sender, 
                    yesVotes: 0, noVotes: 0});
        emit NomineeProposed(_nomineeAddress,  msg.sender);
        emit MonikerAnnounce(_nomineeAddress, _moniker);
    }



    ///@notice Cast a vote for a nominee. Can only be run by people on the whitelist. 
    ///@param _nomineeAddress The address of the nominee
    ///@param _accepted Whether the vote is to accept (true) or reject (false) them. 
    ///@return returns true if the vote has reached a decision, false if not
    ///@return only meaningful if the other return value is true, returns true if the nominee is now on the whitelist. false otherwise.
    function castNomineeVote(address _nomineeAddress, bool _accepted) public payable checkAuthorised(msg.sender) returns (bool, bool){
        
        bool decided = false;
        bool voteresult = false;
        
//      Check if open nominee, other checks redundant
        if (isNominee(_nomineeAddress)) {
            
            
//      Check that this sender has not voted before. Initial config is no redos - so just reject            
            if (address(nomineeList[_nomineeAddress].vote[msg.sender].voter) == address(0)) {
                // Vote is valid. So lets cast the Vote
                nomineeList[_nomineeAddress].vote[msg.sender] = NomineeVote({voter: msg.sender, accept: _accepted });
                
                // Amend Totals
                if (_accepted)
                {
                    nomineeList[_nomineeAddress].yesVotes++;
                } else {    
                    nomineeList[_nomineeAddress].noVotes++;
                }
                
                emitNomineeVoteCastEvent(_nomineeAddress, nomineeList[_nomineeAddress].yesVotes, nomineeList[_nomineeAddress].noVotes, _accepted);

                // Check to see if enough votes have been cast for a decision
                (decided, voteresult) = checkForNomineeVoteDecision(_nomineeAddress);
            }
            else
            {
               emit DebugMessageEvent(_nomineeAddress, ERR_ALREADY_VOTED, address(nomineeList[_nomineeAddress].vote[msg.sender].voter));
            }
        }
        else
        {   // Not a nominee, so set decided to true
            emit DebugMessageEvent(msg.sender, ERR_NOT_ON_NOMINEELIST, _nomineeAddress);
            decided = true;
        }
        
        
        // If decided, check if on whitelist
        if (decided) {
            voteresult = isWhitelisted(_nomineeAddress);
        }
        
        return (decided, voteresult);
        
    }    
    


// These function encapsulates the logic as to whether a vote is complete    



    ///@notice This function encapsulates the logic for percentage completion for the vote and emits the NomineeVoteCast Event. 


    function emitNomineeVoteCastEvent(address _nomineeAddress, uint _yesVotes, uint _noVotes,  bool _accepted) private 
    {
         uint8 progress = uint8((100*_yesVotes)/whiteListCount);
		  	
         emit NomineeVoteCast(_nomineeAddress, msg.sender,_yesVotes, _noVotes, progress, _accepted); 
    }



    ///@notice This function encapsulates the logic for determining if there are enough votes for a definitive decision
    ///@param _nomineeAddress The address of the NomineeElection
    ///@return returns true if the vote has reached a decision, false if not
    ///@return only meaningful if the other return value is true, returns true if the nominee is now on the whitelist. false otherwise.

    function checkForNomineeVoteDecision(address _nomineeAddress) private returns (bool, bool)
    {
        NomineeElection memory election = nomineeList[_nomineeAddress];
        bool decided = false;
        bool voteresult = false;
        
        
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
    

// Adds the user to the white list.


    ///@notice This private function adds the accepted nominee to the whitelist. 
    ///@param _nomineeAddress The address of the nominee being added to the whitelist
    function acceptNominee(address _nomineeAddress) private
    {
       whiteList[_nomineeAddress] = WhitelistPerson(_nomineeAddress, 0);
       whiteListCount++;
       
    // Remove from nominee list   
       removeNominee(_nomineeAddress);
    }
    
    
// Remove person from white list. Not currently used, but will be needed.     
    ///@notice This private function adds the removes a user from the whitelist. Not currently used. 
    ///@param _address The address of the nominee being removed to the whitelist
    
    function deWhiteList(address _address) private 
    {
        
        if(isWhitelisted(_address))
        {
            delete(whiteList[_address]);
            whiteListCount--;
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

    }


}    

