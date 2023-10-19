# Magic Chat

*W.I.P*

Your coworker is working in the server room.  you can barely hear eachother but still need to communicate to ensure you both do what you need to.  Both of your devices are on the same network but despite this you have no good way to talk between them.

Now there is Magic Chat, the local network chatroom inspired by [Magic Wormhole](https://github.com/magic-wormhole/magic-wormhole) Magic Chat uses a 5 word phrase that you come up with to allow you to connect to another user on your network without the address of either machine.  

Simply install the server on your network, modify the client code to point to that server, and compile the program, then Magic Chat is ready to go.


## How It Works

After launching the Magic Chat program, you enter a 5 word key in the format "one-two-three-four-five".
The first word is used to select the channel you wish to talk on, this decides what machine you will attempt to connect with.  The following 3 words are hashed and used as your AES encryption key. To verify that you are both attempting to communicate using the same key, the fith word is appended to your first three, then the hash of that string is taken and sent to the server, along with the channel, and IP:PORT of the machine you are using wormhole from.

If another device tries to connect on the same channel, these hashes are compared.  If they match the port and ip of the other machine is sent to eachother.   These two machines will then attempt to connect to eachother on their own in order to communicate directly.

## Project for KnightHacks 2023