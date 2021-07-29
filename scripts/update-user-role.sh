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

read -p 'Roll No: ' rollno
read -p 'Role Code: ' role_code

sqlite3 iitkcoin.db "UPDATE ACCOUNT SET role=$role_code WHERE rollno=$rollno";
