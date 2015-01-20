# chefstarter
Launch chef-client remotely, with HTTP Basic Auth and TLS support.  Written in Go.

## Why
I run a technical operations team responsible for over 400 Linux instances that are provisioned by Chef. We regularly find ourselves needing to invoke a chef-client run manually across a large number of machines.  Before `chefstarter`, we used [dsh](http://sourceforge.net/projects/dsh/) to SSH into the instances and kick off `chef-client`.  This technique has several disadvanteges:

* It requires SSH keys to be provisioned if you're running it from unattended automation/deployment scripts.  SSH authentication is appropriate for most systems work but for merely launching the `chef-client` process, less intensive and less secure authentication may be perfectly fine.
* It requires a SSH client.  `chefstarter`, on the other hand, receives its commands over HTTP and you only need cURL to kick off a Chef run on a node.
* If you manage large numbers of nodes, your SSH host keys are continually being added/changed, so you have to either pre-add them (painful) or be prepared to accept new keys manually when they come along.

## How
`chefstarter` does one thing and one thing only: it launches `chef-client` when an authenticated HTTP request is made to a pre-arranged URL path.  You run it as a daemon (preferably under something like [supervisord](http://supervisord.org/).  You run it as a non-root user with `sudo` NOPASSWD access to execute `chef-client` (and *only* `chef-client`)

## Command-line Flags
```
./chefstarter -h
Usage of ./chefstarter:
  -listen=":7100": IP:port to listen on (Default: listen on all interfaces on port 7100)
  -wait=false: Wait until chef-client completes run before returning HTTP response
  -user="chefstarter": HTTP Basic Auth username (Default: chefstarter)
  -pass="": HTTP Basic Auth password
  -path="/": Request path to initiate chef-client run (Default: /)
  -ssl=false: Enable HTTPS (true/false)
  -cert="": HTTPS X.509 Public Certificate file
  -key="": HTTPS X.509 Private Key file
```

## Synchronous vs Asynchronous Operation
If you set the `-wait` flag to `true`, chefstarter will execute `chef-client` synchronously and capture the exit code when it finishes.   If there's an error executing `chef-client` or if it returns a non-zero exit code, `chefstarter` will return a 500 HTTP response code.  If `chef-client` executes and returns an exit code of zero (no errors), `chefstarter` returns 200 OK.

## Path
I recommend that you use the `-path` setting to specify a secret path that will trigger the `chef-client` run.  If you set `-path=/startchefnow`, the Chef run will only kick off if a `GET /startchefnow` request is received.  If any other path is requested, it simply returns `404 Not Found`.  The default is to kick off a run when `GET /` is received--you should probably set it to something that only your organization knows.

## Firewalling
I recommend that you use firewall ACLs to restrict access to chefstarter's TCP port.


## Author
`chefstarter` was created by [Chris Snell](http://output.chrissnell.com)
