wampirc
=======


`wampirc` is a [WAMPv2](http://wamp.ws/) <-> IRC bot.

installation
------------

    go install github.com/beatgammit/wampirc

status
======

Very basic functionality works, including:

- connect to a basic `turnpike`-based WAMP server ([tested with example chat server in turnpike](https://github.com/jcelliott/turnpike/tree/v2/examples/chat/chatserver))
- connect to [ngircd](https://github.com/alexbarton/ngircd) with and without a password
- chats go both directions with [turnpike example client](https://github.com/jcelliott/turnpike/tree/v2/examples/chat/chatserver)

There's still quite a lot to do...

license
=======

Licensed under the BSD 3-clause licesne. See LICENSE.BSD3 for details.
