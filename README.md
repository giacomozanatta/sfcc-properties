# SFCC-PROPERTIES

With this simple script you can export in an excel file all the active properties of a Salesforce Commerce Cloud (SFCC) project.  

## Usage

1. First install the program: ```go install github.com/giacomozanatta/sfcc-properties@latest```
2. Navigate in the main folder of your project
3. Create a file config.json. This json is compose of two attributes:
  - cartridge: [string]: you must insert all active cartridge, in reverse overriding order (from the least important to the most one)
  - locales: [string] insert all locales you have in the properties files
 4. Launch the program with ```sfcc-properties```
 5. The program will produce a properties.xlsx file in the main folder of your project.
