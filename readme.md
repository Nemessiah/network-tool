# network-tool

Network Tool is a CLI utility for generating vendor-specific network configuration commands from structured input. It is designed to reduce manual command creation errors by validating input and automatically replacing variables in predefined command templates.

The tool is currently at **version 0.5** and is under active development.

The primary goal of this project is to learn **Go** and software development in general. In addition, it aims to provide a flexible, vendor-agnostic way to generate validated network configuration commands from structured input.

---

## Features (v0.5)

* Command-line interface for generating network commands
* Customizable command schema per vendor
* Validation of custom command definitions
* Automatic replacement of input keys in command templates
* Grouped command output per vendor

---

## Default Supported Platforms

* Palo Alto Networks firewalls
* Cisco switches
* NetBox (planned, not yet implemented)

---

## How It Works

1. Vendor-specific commands are defined using a customizable schema.
2. The user is prompted for required input values.
3. Inputs are validated against command requirements.
4. Keys are replaced in command templates.
5. Commands are output and grouped by vendor.

---

## Example Output

```
commands for: paloalto
    set network interface ethernet1/1 layer3 ip 10.0.0.1/24
    set zone trust network layer3 ethernet1/1

commands for: cisco
    interface GigabitEthernet0/1
    ip address 10.0.0.1 255.255.255.0
```

---

## Roadmap

### Completed

* [x] CLI tool that returns commands (v0.1)
* [x] Customizable command schema/system (v0.2)
* [x] Validation of custom commands (v0.3)
* [x] Automatic replacement of keys in commands (v0.4)
* [x] Program functioning (v0.5)
* [x] Tests for all functions (v0.6)
* [x] All v0 goals complete, and tested (v1.0)

### Planned
* [ ] additional input methods (v1.1)

  * [ ] JSON
  * [ ] CLI parameters
* [ ] Auto-generate input based on customizable standards (v1.2)
* [ ] Query NetBox for existing network and device information (v1.3)
* [ ] Authentication (v2.0-beta)

  * [ ] API token access
* [ ] Write to NetBox (v2.1-beta)
* [ ] RBAC (v2.2-beta)

  * [ ] Read existing
  * [ ] Update commands
* [ ] Expose API and containerize for web deployment (v2.0)
* [ ] GUI frontend (v2.1)
* [ ] Storage of device data (v3.0-beta)
* [ ] Run commands against provided systems automatically (v3.0)

---

## Important Notes

* NetBox integration is planned but **not implemented** as of v0.5.

⚠️ **Warning:** As of v0.5, the **default vendor command definitions have not been fully validated for real-world accuracy**. Generated commands should be reviewed and tested before use in production environments.

---

## Installation

### Requirements

* Go **1.25.4** or newer

### Build

```bash
go build .
```

### Run

```bash
go run .
```

---

## Configuration

Command schemas are defined using **YAML configuration files**. These schemas describe vendors, device types, and ordered command definitions grouped by feature and CRUD action.

Multiple configuration files may be used. If the same vendor name is defined more than once, the **last loaded definition wins**.

---

### Vendor Definition Structure

Each vendor must follow this structure:

```yaml
<vendor_name>:
  devicetype: <type>
  commands:
    - feature: <feature_name>
      actions:
        <action>: [ <command>, <command>, ... ]
```

#### Vendor fields

* **vendor_name** (required)

  * Logical vendor identifier (e.g. `paloalto`, `cisco`, `juniper`)
  * Vendor order does not matter

* **devicetype** (required)
  Must be one of:

  * `firewall`
  * `switch`
  * `api`
  * `dhcp`

* **commands** (required)

  * An **ordered list** of command features
  * Feature order is preserved exactly as written and determines output order

* **feature** (required)

  * Logical grouping of commands (e.g. `interface`, `vlan`, `dhcprelay`)

* **actions** (required)

  * CRUD-style actions supported by this feature
  * Only **one action** is selected per program run

Allowed actions:

* `create`
* `read`
* `update`
* `delete`

Each action must contain a list of **one or more command strings**. Lists must be used even if only a single command exists.

---

### Command Placeholders

Commands may contain placeholder keys that are replaced at runtime using user-provided input.

**Allowed placeholder keys:**

* `interface`
* `ipaddress`
* `vlanid`
* `vlanname`
* `zone`
* `cidr`
* `mask`

Placeholder format:

```text
{{key}}
```

Example:

```text
set interface {{interface}} ip {{ipaddress}}
```

---

### Validation Behavior

* Invalid CRUD action names are ignored with a warning
* Invalid placeholder keys cause **configuration validation to fail**
* Features with no valid actions are ignored

---

### Current Prompting Behavior (v0.5)

* Interactive mode only
* The program prompts for **every allowed placeholder key** on each run
* After selecting a CRUD action, **all commands** for that action are output per vendor
