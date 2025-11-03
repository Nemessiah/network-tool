Tool for getting command output for a few different tools. Mainly Palo Alto Firewalls, Cisco switches, and Netbox.

> Currently in v0.0
{.is-warning}

Required inputs

required:
- General
    - VLAN
    - Managing Site
    - Site
    - CIDR
- Program
    - Network type
        - Guest
        - Client
        - Wan
        - Transit
        - DMZ
        - user
- Switch
    - trunk ports
    - Access ports
- Firewall
    - interface name
    - vrf
    - user id
- DHCP
    - range
    - lease time
- Netbox
    - existing IDs