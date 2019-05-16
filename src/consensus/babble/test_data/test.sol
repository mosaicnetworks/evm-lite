pragma solidity 0.5.7;

contract Test {

  function checkAuthorisedPublicKey(bytes32 _publicKey) public pure returns (bool) {
      
      if(_publicKey == "0x12345") {
          return true;
      } else {
          return false;
      }
   }
}