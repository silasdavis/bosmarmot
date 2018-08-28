pragma solidity ^0.4.24;

contract simplestorage {
    uint public storedData;

    constructor(uint initVal) public {
        storedData = initVal;
    }

    function set(uint value) public {
        storedData = value;
    }

    function get() public constant returns (uint) {
        return storedData;
    }

    // Since transactions are executed atomically we can implement this concurrency primitive in Solidity with the
    // desired behaviour
    function testAndSet(uint expected, uint value) public returns (uint, bool) {
        if (storedData == expected) {
            storedData = value;
            return (storedData, true);
        }
        return (storedData, false);
    }
}