#!/bin/bash

echo "
UPDATE USER ROLE

Roll Codes
========================
0 NormalUser       
1 GeneralSecretary 
2 AssociateHead    
3 CoreTeamMember   
"
read -p 'Role Code: ' role_code

echo "
The following roll no.s have the role code $role_code:
"
sqlite3 iitkcoin.db "SELECT rollno FROM ACCOUNT WHERE role=$role_code"