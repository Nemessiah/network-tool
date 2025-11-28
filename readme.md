Tool for getting command output for a few different tools. Mainly Palo Alto Firewalls, Cisco switches, and Netbox.

> Currently in v0.0
{.is-warning}

Required inputs

required:
- General
    - VLAN
    - Site
    - CIDR
- Program
    - Network type
        - Guest
        - Client
        - Wan
        - Transit
        - DMZ
        - server
        - device
- Switch
    - trunk ports
    - Access ports
- Firewall
    - interface name
    - interface settings
    - vrf
    - user id
    - DHCP Relay
    - 
- DHCP
    - range
    - lease time
- Netbox
    - existing IDs
    - Add new
        - networks
        - interfaces
Nice two have:

