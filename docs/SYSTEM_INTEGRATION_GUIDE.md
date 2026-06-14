# NAS Dashboard - System Integration & Refinement

## Introduction
This document details the refined interconnected logic of the NAS Dashboard. Every module has been integrated with "Safety Interlocks" and "Cross-Service Synchronization" to ensure a robust, production-ready NAS environment.

## 1. Storage & File Sharing Chain
The primary path from hardware to user access is fully automated and protected:
- **Physical Disk -> Storage Pool**: New disks are automatically formatted (ext4) and mounted to managed paths (`/mnt/storage_pools/...`) before being merged into the unified MergerFS view.
- **Persistence**: Every MergerFS pool generates a native `systemd` unit (`nas-mergerfs-[name].service`) to guarantee boot-time availability.
- **Safety Interlock**: The system prevents deleting a Storage Pool if any SMB share or Sync job is still referencing its mount point.

## 2. Identity & Security Integration
- **Unified Identity**: When a NAS user is created via the UI, the system simultaneously executes `useradd` (Linux level) and `smbpasswd` (Samba level). This ensures that login credentials remain synchronized across the system and network shares.
- **Cleanup**: Deleting a user automatically removes their Samba credentials, maintaining a clean security posture.

## 3. Network Discovery & macOS Integration
- **Time Machine (Auto-mDNS)**: Mark any SMB share as a "Time Machine" target, and the system handles the heavy lifting:
  - Injects specific AFP/Apple fruit parameters into `smb.conf`.
  - Automatically generates an Avahi mDNS service XML file (`/etc/avahi/services/...`).
  - Restarts the Avahi daemon to make the NAS instantly discoverable by macOS Finder.

## 4. Data Protection Ecosystem
- **Sync (Rsync)**: Managed through a background job system that stores task states in SQLite. It targets the persistent mount points created by the storage module.
- **Backup (Restic)**: Provides deduplicated snapshots to local or cloud repositories.

## Summary of Refinements
| Integration Point | Logic Applied | Benefit |
| :--- | :--- | :--- |
| **User Management** | Linux User ↔ Samba Sync | Zero-config SMB access after user creation. |
| **Storage Deletion** | Check active SMB/Sync Jobs | Prevents accidental mount-point breakage. |
| **Time Machine** | Samba Config + Avahi Broadcast | Native Apple device compatibility. |
| **MergerFS** | Auto-Format + systemd Persistence | Raw disks to boot-persistent pools in one click. |

## Conclusion
The system is now a cohesive unit. The "glue" code implemented ensures that high-level UI actions correctly ripple down through all necessary Linux system services, maintaining state consistency and data safety.
