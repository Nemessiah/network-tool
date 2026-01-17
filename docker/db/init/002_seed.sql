-- ------------------------------------------------------------------
-- Vendors
-- ------------------------------------------------------------------
INSERT INTO vendors (name, devicetype)
VALUES
  ('paloalto', 'firewall'),
  ('cisco', 'switch');

-- ------------------------------------------------------------------
-- Features
-- ------------------------------------------------------------------

-- Palo Alto
INSERT INTO features (vendor_id, name)
SELECT id, 'interface' FROM vendors WHERE name = 'paloalto';

INSERT INTO features (vendor_id, name)
SELECT id, 'dhcprelay' FROM vendors WHERE name = 'paloalto';

-- Cisco
INSERT INTO features (vendor_id, name)
SELECT id, 'vlan' FROM vendors WHERE name = 'cisco';

-- ------------------------------------------------------------------
-- Actions
-- ------------------------------------------------------------------

-- Palo Alto: interface → create
INSERT INTO actions (feature_id, action)
SELECT f.id, 'create'
FROM features f
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface';

-- Palo Alto: interface → update
INSERT INTO actions (feature_id, action)
SELECT f.id, 'update'
FROM features f
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface';

-- Palo Alto: dhcprelay → create
INSERT INTO actions (feature_id, action)
SELECT f.id, 'create'
FROM features f
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'dhcprelay';

-- Cisco: vlan → create/read/update/delete
INSERT INTO actions (feature_id, action)
SELECT f.id, a
FROM features f
JOIN vendors v ON v.id = f.vendor_id
CROSS JOIN UNNEST(ARRAY['create','read','update','delete']::action_type[]) a
WHERE v.name = 'cisco' AND f.name = 'vlan';

-- ------------------------------------------------------------------
-- Initial action revisions + set current_revision_id
-- ------------------------------------------------------------------

-- Create one initial revision per action
INSERT INTO action_revisions (action_id, created_by, comment)
SELECT id, 'seed', 'initial revision'
FROM actions;

-- Point each action at its revision
UPDATE actions a
SET current_revision_id = ar.id
FROM action_revisions ar
WHERE ar.action_id = a.id;

-- ------------------------------------------------------------------
-- Commands (insert against the action's current_revision_id)
-- ------------------------------------------------------------------

-- Palo Alto: interface → create
INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 1,
       'set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'create';

INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 2,
       'set network interface aggregate-ethernet {{interface}} description {{vlanname}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'create';

INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 3,
       'set zone {{zone}} network layer3 {{interface}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'create';

-- Palo Alto: interface → update (same commands)
INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 1,
       'set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'update';

INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 2,
       'set network interface aggregate-ethernet {{interface}} description {{vlanname}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'update';

INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 3,
       'set zone {{zone}} network layer3 {{interface}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name = 'paloalto' AND f.name = 'interface' AND a.action = 'update';

-- Cisco: vlan → create
INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 1, 'vlan {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='create';

INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 2, 'name {{vlanname}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='create';

-- Cisco: vlan → read
INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 1, 'show vlan id {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='read';

-- Cisco: vlan → update
INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 1, 'vlan {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='update';

INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 2, 'name {{vlanname}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='update';

-- Cisco: vlan → delete
INSERT INTO commands (revision_id, position, command)
SELECT a.current_revision_id, 1, 'no vlan {{vlanid}}'
FROM actions a
JOIN features f ON f.id = a.feature_id
JOIN vendors v ON v.id = f.vendor_id
WHERE v.name='cisco' AND f.name='vlan' AND a.action='delete';
