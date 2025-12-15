Tool for getting command output for a few different tools. Mainly Palo Alto Firewalls, Cisco switches, and Netbox.

> [!WARNING]
> Currently in v0.3 and not functional

Required inputs:
 - interface 
 - ipaddress
 - vlanid
 - vlanname
 - zone
 - cidr
 - mask

Current Goals:
- [x] CLI tools that returns commands > 0.1
- [x] Customizable commands schema/system > 0.2
- [x] Validation of custom commands > 0.3
- [x] Automatic replacement of key's in commands > 0.4
- [] Tests for all functions > 0.5
- [] Generate commands via CLI input. > 1.0
- [] Auto generate input based on internal standards. > 1.1
- [] Query Netbox for existing network, and device information. > 1.3
- [] Authentication. 2.0-beta
    - [] API token access
- [] write to netbox 2.1-beta
- [] RBAC 2.2-beta
    - [] read existing
    - [] update commands
- [] Expose API and containerize for web deployment. > 2.0
- [] GUI front end. > 2.1
- [] storage of device data 3.0-beta
- [] Run commands agains provided systems automatically. > 3.0

