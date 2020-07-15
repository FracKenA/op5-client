Setup Steps:
# Create group with appropriate rights.

Group name: satellite

OP5 Monitor Group Rights:
    Nagios Auth:
        - Configuration Information
    API:
        - API Command
        - API Config
    Host:
        - Host Add Delete
        - Host View All
        - Host Edit All
    Host Commands:
        - Host Command Passive Check
    Host Template:
        - Host Template View All
    Service:
        - Service Add Delete
        - Service View All
        - Service Edit All
    Service Commands:
        - Service Command Passive Check
    Service Template:
        - Service Template View All
    Configuration:
        - Export

# Create local user.

Username: name
Realname: Satellite Account
Password: pass
Auth. Modules: Default
Groups: satellite

# Create passive host template.

! The check command could vary based on the usage.

Name: passive-host-template
Check Command: check_dummy
Check Command Args: 2!'Host is in a critical state'!'output=2`
Contact Groups: support-group
Check Freshness: Yes
Fresshness Threshold: 0
Active Checks Enabled: No
Passive Checks Enabled: Yes
Flap Detection Enabled: No

# Create passive service template

! This could vary based on the check.

Name: passive-service-template
Check Command: check_dummy
Check Command Args: 2!'Service is in a critical state'!''
Active Checks Enabled: No
Passive Checks Enabled: Yes
Check Freshness: Yes
Freshness Threshold: 0
Flap Detection Enabled: No
Contact Groups: support-group

# Register the server with Monitor.

`satellite register`

