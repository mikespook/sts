# RPCD

This service provides two RPC calls for authorization of `stsd`.

## Auth.Password

This method is used for plain password authorization.

### Input

Type: \*model.Auth

### Reply

Type: \*ssh.Permissions

## Auth.PublicKey

This method is used for public key authorization.

### Input

Type: \*model.Auth

### Reply

Type: \*ssh.Permissions

# STSD

## Ctrl.Cutoff

This method is used for cutting a connection off.

## Ctrl.Kickoff

This method is used for Kicking user off line.

## Ctrl.Restart

This method is used for restart the internal tunnel service.

## Stat.User

This method returns status of a user tunnel, including activity connections, online time and network throughputs.

## Stat.Tunnel

This method returns aggregation status of tunnels, including online users, activity connections and network throughputs.

## Stat.Server

This method returns status of the service, including PID, established time and the last internal error.
