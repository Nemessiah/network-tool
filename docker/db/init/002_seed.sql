-- ------------------------------------------------------------------
-- Vendors
-- ------------------------------------------------------------------

INSERT INTO vendors (name, devicetype)
VALUES
  ('paloalto', 'firewall'),
  ('cisco', 'switch');

-- ------------------------------------------------------------------
-- Palo Alto features
-- ------------------------------------------------------------------

-- interface
INSERT INTO features (vendor_id, name, position)
SELECT id, 'interface', 1 FROM vendors WHERE name = 'paloalto';

-- dhcprelay
INSERT INTO features (vendor_id, name, position)
SELECT id, 'dhcprelay', 2 FROM vendors WHERE name = 'paloalto';

-- ------------------------------------------------------------------
-- Palo Alto actions
-- ------------------------------------------------------------------

-- interface → create
INSERT INTO actions (feature_id, action)
SELECT f.id, 'create'
FROM features f
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface';

-- interface → update
INSERT INTO actions (feature_id, action)
SELECT f.id, 'update'
FROM features f
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface';

-- dhcprelay → create
INSERT INTO actions (feature_id, action)
SELECT f.id, 'create'
FROM features f
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'dhcprelay';

-- ------------------------------------------------------------------
-- Palo Alto commands
-- ------------------------------------------------------------------

-- interface → create
INSERT INTO commands (action_id, position, command)
SELECT a.id, 1,
       'set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'create';

INSERT INTO commands (action_id, position, command)
SELECT a.id, 2,
       'set network interface aggregate-ethernet {{interface}} description {{vlanname}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'create';

INSERT INTO commands (action_id, position, command)
SELECT a.id, 3,
       'set zone {{zone}} network layer3 {{interface}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'create';

-- interface → update (same commands)
INSERT INTO commands (action_id, position, command)
SELECT a.id, 1,
       'set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'update';

INSERT INTO commands (action_id, position, command)
SELECT a.id, 2,
       'set network interface aggregate-ethernet {{interface}} description {{vlanname}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'update';

INSERT INTO commands (action_id, position, command)
SELECT a.id, 3,
       'set zone {{zone}} network layer3 {{interface}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'update';

-- ------------------------------------------------------------------
-- Cisco features
-- ------------------------------------------------------------------

INSERT INTO features (vendor_id, name, position)
SELECT id, 'vlan', 1 FROM vendors WHERE name = 'cisco';

-- ------------------------------------------------------------------
-- Cisco actions
-- ------------------------------------------------------------------

INSERT INTO actions (feature_id, action)
SELECT f.id, a
FROM features f
JOIN vendors v ON v.id = f.vendor_id
CROSS JOIN UNNEST(ARRAY['create','read','update','delete']::action_type[]) a
WHERE v.name = 'cisco' AND f.name = 'vlan';

-- ------------------------------------------------------------------
-- Cisco commands
-- ------------------------------------------------------------------

-- create
INSERT INTO commands (action_id, position, command)
SELECT a.id, 1, 'vlan {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='create';

INSERT INTO commands (action_id, position, command)
SELECT a.id, 2, 'name {{vlanname}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='create';

-- read
INSERT INTO commands (action_id, position, command)
SELECT a.id, 1, 'show vlan id {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='read';

-- update
INSERT INTO commands (action_id, position, command)
SELECT a.id, 1, 'vlan {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='update';

INSERT INTO commands (action_id, position, command)
SELECT a.id, 2, 'name {{vlanname}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='update';

-- delete
INSERT INTO commands (action_id, position, command)
SELECT a.id, 1, 'no vlan {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='delete';
