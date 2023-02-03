# SMB-Protect 

## Why
Consider a scenario. One day your Windows Machine is infected with a Ransomware.
This one encrypts all your files and asks you to wire 100 Bitcoins to $ADDR for 
getting it decrypted. 

You rejoice at your insight, for, you backed up all your data to a NAS...
...the NAS that your computer had cached credentials for...

You hurry and logon to NAS to see disaster. All the data has been taken hostage 
there as well. You despair.

If only you knew about this nifty trick here.

## How it works
Create 2 users on your share. Grant one of them read-only permissions and the other 
read-write.

SMB-Protect will start at login and mount all your shares read-only.
As you usually don't edit the files on the drives, there is no need for writing to be enabled.

When you do need to update files, you can click a button on the systray to re-map the drive as 
a writable mount. When you are done with your work, you can switch it back to Read-only