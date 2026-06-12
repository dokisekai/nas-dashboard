# NAS Dashboard - User Guide

Complete user guide for the NAS Dashboard, covering all features, functionality, and daily operations.

## 📋 Table of Contents

1. [Getting Started](#getting-started)
2. [Desktop Environment](#desktop-environment)
3. [System Monitoring](#system-monitoring)
4. [Storage Management](#storage-management)
5. [User Management](#user-management)
6. [Application Center](#application-center)
7. [Plugin System](#plugin-system)
8. [Settings and Configuration](#settings-and-configuration)
9. [Tips and Tricks](#tips-and-tricks)
10. [Troubleshooting](#troubleshooting)

---

## Getting Started

### First Login

#### Accessing the Dashboard

1. **Open your web browser**
2. **Navigate to your dashboard URL**
   - Local: `http://localhost:3000`
   - Remote: `https://your-domain.com`

3. **Login with default credentials**
   - Username: `admin`
   - Password: `admin123`

4. **Change your password immediately**
   - Click on user icon in top right
   - Select "Change Password"
   - Enter new secure password

### Dashboard Tour

#### Main Interface

```
┌─────────────────────────────────────────────────────────────┐
│  [Logo]  NAS Dashboard                    [🔔] [👤] [⚙️]     │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│                    Desktop Area                               │
│                                                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Widget    │  │   Widget    │  │   Widget    │         │
│  │             │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Widget    │  │   Widget    │  │   Widget    │         │
│  │             │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                               │
├─────────────────────────────────────────────────────────────┤
│  [🏠] [💾] [📊] [👥] [⚙️]                                    │
│  Home  Storage  Monitor  Users  Settings                     │
└─────────────────────────────────────────────────────────────┘
```

#### Key Components

- **Top Bar**: Logo, notifications, user menu, settings
- **Desktop Area**: Main workspace with widgets and windows
- **Dock Bar**: Quick access to applications
- **Context Menu**: Right-click for additional options

---

## Desktop Environment

### Desktop Features

#### Window Management

**Opening Windows:**

1. Click on any application icon in the dock
2. Window opens on desktop
3. Multiple windows can be open simultaneously

**Window Controls:**

- **Title Bar**: Drag to move window
- **Minimize**: `-` button (minimizes to dock)
- **Maximize**: `□` button (fullscreen)
- **Close**: `×` button (closes window)

**Window Snapping:**

Windows automatically snap to edges when dragged:

```
┌────────────────────────────────────┐
│          Left Half                 │
│                                    │
└────────────────────────────────────┘
┌────────────────────────────────────┐
│          Right Half                │
│                                    │
└────────────────────────────────────┘
```

**Keyboard Shortcuts:**

- `Win + ↑`: Maximize window
- `Win + ↓`: Restore/Minimize window
- `Win + ←`: Snap to left half
- `Win + →`: Snap to right half
- `Alt + F4`: Close window

#### Widget System

**Adding Widgets:**

1. **Right-click on desktop**
2. **Select "Add Widget"**
3. **Choose widget from library**
4. **Configure widget settings**
5. **Click "Add"**

**Available Widgets:**

- **System Monitor**: CPU, memory, disk usage
- **Weather**: Current weather conditions
- **Calendar**: Upcoming events
- **Clock**: Time and date display
- **Quick Note**: Simple notepad
- **Storage**: Disk space overview
- **Network**: Network statistics

**Widget Operations:**

- **Move**: Drag widget to new position
- **Resize**: Drag bottom-right corner
- **Configure**: Click gear icon
- **Remove**: Click `×` icon

#### Desktop Customization

**Change Background:**

1. Right-click on desktop
2. Select "Change Background"
3. Choose from:
   - Solid colors
   - Gradients
   - Custom image

**Desktop Settings:**

1. Right-click on desktop
2. Select "Desktop Settings"
3. Configure:
   - Widget snap to grid
   - Show desktop icons
   - Auto-arrange widgets

#### Dock Bar

**Dock Features:**

```
┌─────────────────────────────────────────────────────────────┐
│  [🏠] [💾] [📊] [👥] [⚙️]                                │
│  Home  Storage  Monitor  Users  Settings                     │
│                                                               │
│  Running Applications: [📁] [📊]                            │
└─────────────────────────────────────────────────────────────┘
```

**Dock Operations:**

- **Launch App**: Click icon
- **Pin/Unpin**: Right-click > "Pin to Dock"
- **Reorder**: Drag icon to new position
- **Remove**: Drag icon off dock

---

## System Monitoring

### Real-time Monitoring

#### CPU Monitoring

**View CPU Usage:**

1. Open System Monitor app
2. Navigate to CPU tab
3. View real-time metrics:
   - Overall CPU usage
   - Per-core usage
   - Load averages (1, 5, 15 min)
   - CPU frequency
   - Temperature

**CPU Chart:**

```
CPU Usage
┌─────────────────────────────────────────────────────────────┐
│  ████████████████████████████████████████████████████  85% │
│                                                             │
│  Core 1: ████████████████████████████████████████████  80% │
│  Core 2: ████████████████████████████████████████████  82% │
│  Core 3: ████████████████████████████████████████████  78% │
│  Core 4: ██████████████████████████████████████████  75%  │
└─────────────────────────────────────────────────────────────┘
```

#### Memory Monitoring

**View Memory Usage:**

1. Open System Monitor app
2. Navigate to Memory tab
3. View memory statistics:
   - Total memory
   - Used memory
   - Free memory
   - Swap usage
   - Memory percentage

**Memory Display:**

```
Memory Usage
┌─────────────────────────────────────────────────────────────┐
│  Total: 16.0 GB  Used: 8.2 GB  Free: 7.8 GB                │
│                                                             │
│  ████████████████████████████████████████░░░░░░░░░░  51%   │
│                                                             │
│  Swap: 2.0 GB  Used: 0.0 GB  Free: 2.0 GB                  │
└─────────────────────────────────────────────────────────────┘
```

#### Disk Monitoring

**View Disk Usage:**

1. Open System Monitor app
2. Navigate to Disk tab
3. View disk information:
   - Disk partitions
   - Usage per partition
   - File system type
   - Mount points
   - I/O statistics

**Disk Display:**

```
Disk Usage
┌─────────────────────────────────────────────────────────────┐
│  / (root)                                                    │
│  Total: 100 GB  Used: 45 GB  Free: 55 GB                    │
│  ████████████████████████░░░░░░░░░░░░░░░░░░░░░░  45%       │
│                                                             │
│  /home                                                       │
│  Total: 500 GB  Used: 200 GB  Free: 300 GB                 │
│  ████████████████████████████████████████░░░░░░░░░  40%   │
└─────────────────────────────────────────────────────────────┘
```

#### Network Monitoring

**View Network Statistics:**

1. Open System Monitor app
2. Navigate to Network tab
3. View network data:
   - Interface status
   - Data transferred (sent/received)
   - Packets transferred
   - Errors and drops
   - Connection speed

**Network Display:**

```
Network Activity
┌─────────────────────────────────────────────────────────────┐
│  eth0: Connected                                             │
│                                                             │
│  Upload:   ████████░░░░░░░░░░░  1.2 MB/s                   │
│  Download: ████████████████████████████░░  8.5 MB/s       │
│                                                             │
│  Total Sent:     125.5 GB                                   │
│  Total Received: 1.2 TB                                     │
└─────────────────────────────────────────────────────────────┘
```

### Historical Data

**View Historical Trends:**

1. Open System Monitor app
2. Select time range:
   - Last hour
   - Last 6 hours
   - Last 24 hours
   - Last 7 days
3. View historical charts
4. Export data to CSV

---

## Storage Management

### Disk Management

#### View Disks

**List All Disks:**

1. Open Storage Manager app
2. Navigate to "Disks" tab
3. View disk information:
   - Disk name and size
   - Model and serial number
   - Partitions
   - Health status

**Disk Display:**

```
Disks
┌─────────────────────────────────────────────────────────────┐
│  /dev/sda (500 GB)                                          │
│  Model: Samsung SSD 860 EVO                                │
│  Status: Healthy                                            │
│                                                             │
│  Partition    Size    Type    Used    Mount Point           │
│  /dev/sda1    100 GB  ext4    45 GB   /                    │
│  /dev/sda2    400 GB  ext4    200 GB  /home                │
└─────────────────────────────────────────────────────────────┘
```

#### Mount/Unmount Disks

**Mount Disk:**

1. Select unmounted partition
2. Click "Mount" button
3. Choose mount point
4. Select mount options
5. Click "Mount"

**Unmount Disk:**

1. Select mounted partition
2. Click "Unmount" button
3. Confirm unmount

#### Format Disk

**Format Partition:**

1. Select partition
2. Click "Format" button
3. Choose file system:
   - ext4 (recommended for Linux)
   - NTFS (Windows compatibility)
   - exFAT (external drives)
4. Enter label
5. Confirm format (⚠️ destroys all data)

### SMB Shares

#### Create SMB Share

**Add New Share:**

1. Open Storage Manager app
2. Navigate to "SMB Shares" tab
3. Click "Add Share"
4. Configure share:
   - Share name
   - Path to share
   - Description
   - Access permissions
   - Guest access
5. Click "Create"

**Share Configuration:**

```
SMB Share Configuration
┌─────────────────────────────────────────────────────────────┐
│  Share Name: documents                                       │
│  Path: /home/shared/documents                               │
│  Description: Shared documents folder                       │
│                                                             │
│  Permissions:                                               │
│  ☑ Read All Users                                           │
│  ☐ Write All Users                                          │
│  ☐ Guest Access                                             │
│                                                             │
│  [Cancel]  [Create]                                         │
└─────────────────────────────────────────────────────────────┘
```

#### Manage SMB Shares

**Edit Share:**

1. Select share from list
2. Click "Edit" button
3. Modify settings
4. Click "Save"

**Delete Share:**

1. Select share from list
2. Click "Delete" button
3. Confirm deletion

### File Browser

**Browse Files:**

1. Open Storage Manager app
2. Navigate to "Files" tab
3. Browse file system
4. View file details:
   - Name and size
   - Type and permissions
   - Owner and group
   - Modification date

**File Operations:**

- **Navigate**: Click folders to browse
- **Upload**: Drag files or click "Upload"
- **Download**: Select file > Click "Download"
- **Delete**: Select file > Click "Delete"
- **Rename**: Select file > Click "Rename"

---

## User Management

### User Accounts

#### Create User

**Add New User:**

1. Open User Manager app
2. Click "Add User" button
3. Fill in user details:
   - Username (required)
   - Email (required)
   - Password (required)
   - Confirm password
   - Role (admin/user)
4. Click "Create"

**User Form:**

```
Create User
┌─────────────────────────────────────────────────────────────┐
│  Username:     [_________________]                          │
│  Email:        [_________________]                          │
│  Password:     [_________________]                          │
│  Confirm:      [_________________]                          │
│  Role:         [Admin ▼]                                    │
│                                                             │
│  [Cancel]  [Create]                                        │
└─────────────────────────────────────────────────────────────┘
```

#### Manage Users

**Edit User:**

1. Select user from list
2. Click "Edit" button
3. Modify user details
4. Click "Save"

**Delete User:**

1. Select user from list
2. Click "Delete" button
3. Confirm deletion

**Change Password:**

1. Select user from list
2. Click "Change Password" button
3. Enter new password
4. Confirm password
5. Click "Change"

#### User Roles

**Admin Permissions:**

- Full system access
- User management
- System settings
- Storage management
- Docker management

**User Permissions:**

- View system monitoring
- Access personal files
- Change own password
- View logs

### SSH Keys

#### Add SSH Key

**Upload Public Key:**

1. Select user from list
2. Navigate to "SSH Keys" tab
3. Click "Add SSH Key"
4. Enter key details:
   - Key name
   - Public key (paste)
5. Click "Add"

**SSH Key Form:**

```
Add SSH Key
┌─────────────────────────────────────────────────────────────┐
│  Key Name:     [_________________]                          │
│  Public Key:                                               │
│  [ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...]                  │
│                                                             │
│  [Cancel]  [Add]                                           │
└─────────────────────────────────────────────────────────────┘
```

#### Manage SSH Keys

**View Keys:**

- Key name
- Fingerprint
- Added date

**Delete Key:**

1. Select key
2. Click "Delete" button
3. Confirm deletion

---

## Application Center

### Storage Manager

**Features:**

- Disk management
- SMB share configuration
- File browser
- Storage statistics

**Usage:**

1. Click Storage icon in dock
2. Navigate through tabs:
   - Disks: Manage disk partitions
   - SMB: Configure file sharing
   - Files: Browse file system

### System Monitor

**Features:**

- Real-time CPU monitoring
- Memory usage tracking
- Disk activity display
- Network statistics
- Historical data viewing

**Usage:**

1. Click Monitor icon in dock
2. Select monitoring category:
   - CPU: View processor usage
   - Memory: Check RAM utilization
   - Disk: Monitor disk activity
   - Network: View network stats

### User Manager

**Features:**

- Create/delete users
- Manage user permissions
- SSH key management
- Password management

**Usage:**

1. Click Users icon in dock
2. Manage user accounts
3. Configure SSH access

### System Settings

**Features:**

- Network configuration
- Service management
- System updates
- Backup management

**Usage:**

1. Click Settings icon in dock
2. Navigate settings categories
3. Modify system configuration

### App Center

**Browse Applications:**

1. Open App Center
2. Browse categories:
   - System Tools
   - Utilities
   - Monitoring
   - Productivity
3. View app details
4. Install applications

**Install Application:**

1. Select application
2. Click "Install" button
3. Wait for installation
4. Application appears in dock

---

## Plugin System

### Plugin Store

**Browse Plugins:**

1. Open Plugin Store app
2. Browse available plugins
3. View plugin details:
   - Description
   - Features
   - Screenshots
   - Reviews
   - Version

**Install Plugin:**

1. Select plugin
2. Click "Install" button
3. Review permissions
4. Confirm installation
5. Plugin activates automatically

**Manage Plugins:**

1. Navigate to "My Plugins" tab
2. View installed plugins
3. Enable/disable plugins
4. Configure plugin settings
5. Uninstall plugins

### Plugin Permissions

**Common Permissions:**

- **Storage**: Access plugin storage
- **UI**: Register widgets and apps
- **Network**: Make HTTP requests
- **WebSocket**: Real-time data access
- **System**: System information access

**Review Permissions:**

Always review plugin permissions before installation to ensure security.

---

## Settings and Configuration

### System Settings

#### Network Configuration

**Configure Network:**

1. Open System Settings
2. Navigate to "Network" tab
3. Configure interfaces:
   - IP address
   - Subnet mask
   - Gateway
   - DNS servers

**Network Display:**

```
Network Configuration
┌─────────────────────────────────────────────────────────────┐
│  Interface: eth0                                             │
│  IP Address: 192.168.1.100                                   │
│  Subnet Mask: 255.255.255.0                                  │
│  Gateway: 192.168.1.1                                        │
│  DNS: 8.8.8.8, 8.8.4.4                                       │
│                                                             │
│  [Apply Changes]                                            │
└─────────────────────────────────────────────────────────────┘
```

#### Service Management

**Manage Services:**

1. Open System Settings
2. Navigate to "Services" tab
3. View service list
4. Control services:
   - Start service
   - Stop service
   - Restart service
   - Enable/disable autostart

**Service Display:**

```
System Services
┌─────────────────────────────────────────────────────────────┐
│  Service         Status    State    Actions                  │
│  nginx           running    active   [Restart] [Stop]        │
│  postgresql      running    active   [Restart] [Stop]        │
│  docker          running    active   [Restart] [Stop]        │
│  apache          stopped    inactive [Start]  [Enable]       │
└─────────────────────────────────────────────────────────────┘
```

#### System Updates

**Check for Updates:**

1. Open System Settings
2. Navigate to "Updates" tab
3. Click "Check for Updates"
4. View available updates
5. Click "Install Updates"

**Update Display:**

```
System Updates
┌─────────────────────────────────────────────────────────────┐
│  Current Version: 0.1.0                                     │
│  Latest Version: 0.1.1                                      │
│                                                             │
│  Available Updates:                                         │
│  ☑ nas-dashboard: Update to version 0.1.1                   │
│  ☐ system packages: 5 updates available                     │
│                                                             │
│  [Update All]  [Check for Updates]                         │
└─────────────────────────────────────────────────────────────┘
```

### Backup Management

#### Create Backup

**Manual Backup:**

1. Open System Settings
2. Navigate to "Backup" tab
3. Click "Create Backup"
4. Select backup type:
   - Full backup
   - Database only
   - Configuration only
5. Click "Start Backup"

**Backup Display:**

```
Backup Management
┌─────────────────────────────────────────────────────────────┐
│  Create Backup                                              │
│                                                             │
│  Backup Type: [Full Backup ▼]                               │
│  Description: [_________________]                          │
│                                                             │
│  [Start Backup]                                            │
└─────────────────────────────────────────────────────────────┘
```

#### Restore Backup

**Restore from Backup:**

1. Open System Settings
2. Navigate to "Backup" tab
3. Select backup from list
4. Click "Restore" button
5. Confirm restore
6. System restarts

#### Schedule Backups

**Configure Automatic Backups:**

1. Open System Settings
2. Navigate to "Backup" tab
3. Click "Schedule" button
4. Configure schedule:
   - Frequency (daily, weekly)
   - Time
   - Backup type
   - Retention policy
5. Click "Save Schedule"

---

## Tips and Tricks

### Keyboard Shortcuts

#### System Shortcuts

- `Ctrl + L`: Focus address bar
- `Ctrl + R`: Refresh current view
- `Ctrl + F`: Find in current view
- `Ctrl + ,`: Open settings
- `Ctrl + H`: Open history

#### Window Shortcuts

- `Alt + Tab`: Switch between windows
- `Win + D`: Show desktop
- `Win + L': Lock screen
- `Win + E`: Open file manager

### Desktop Tips

**Widget Organization:**

1. Group related widgets together
2. Place most-used widgets in top-left
3. Use snap-to-grid for alignment
4. Keep widgets organized by category

**Window Management:**

1. Use snapping to organize windows
2. Minimize unused windows
3. Close windows when done
4. Use multiple desktops for organization

### Performance Tips

**Improve Performance:**

1. Close unused applications
2. Disable unnecessary widgets
3. Clear browser cache regularly
4. Use wired network connection
5. Reduce chart update frequency

**Monitor Performance:**

1. Keep System Monitor open
2. Check CPU usage regularly
3. Monitor memory usage
4. Watch disk space
5. Track network activity

### Security Tips

**Best Practices:**

1. Use strong passwords
2. Enable automatic updates
3. Review plugin permissions
4. Regular security audits
5. Backup important data

**Access Control:**

1. Create separate user accounts
2. Limit admin access
3. Use SSH keys instead of passwords
4. Enable two-factor authentication (when available)
5. Regular password changes

---

## Troubleshooting

### Common Issues

#### Login Problems

**Can't Login:**

1. Check username and password
2. Verify caps lock is off
3. Reset password if needed
4. Clear browser cache
5. Try different browser

#### Widget Issues

**Widget Not Loading:**

1. Check network connection
2. Refresh the page
3. Re-add the widget
4. Check system status
5. Review widget logs

#### Application Issues

**App Won't Open:**

1. Check if app is installed
2. Restart the application
3. Clear app cache
4. Check for updates
5. Reinstall application

#### Performance Issues

**Slow Response:**

1. Check system resources
2. Close unused applications
3. Clear browser cache
4. Check network speed
5. Restart system

### Getting Help

**Support Resources:**

1. **Documentation**: Check user guide
2. **Community Forums**: Ask questions
3. **Issue Tracker**: Report bugs
4. **Email Support**: Contact support team

**Debug Mode:**

Enable debug mode for detailed logging:

1. Click user icon
2. Select "Settings"
3. Enable "Debug Mode"
4. View console logs

---

## Conclusion

This user guide covers all major features of the NAS Dashboard. For additional help:

- Check the online documentation
- Visit the community forums
- Contact support team
- Review troubleshooting section

### Quick Reference

**Common Tasks:**

- **Add Widget**: Right-click desktop > Add Widget
- **Open App**: Click icon in dock
- **Manage Users**: Users app > Add/Edit/Delete
- **View Monitor**: Monitor app > Select category
- **Manage Storage**: Storage app > Disks/SMB/Files
- **Install Plugin**: Plugin Store > Select > Install

**Keyboard Shortcuts:**

- `Ctrl + R`: Refresh
- `Alt + Tab`: Switch windows
- `Win + D`: Show desktop
- `Ctrl + ,`: Settings

**Essential URLs:**

- Dashboard: `http://localhost:3000`
- API: `http://localhost:8888/api`
- Documentation: `https://docs.nas-dashboard.com`

---

**Last Updated**: 2026-06-12  
**Version**: 0.1.0  
**Status**: User Guide Complete
