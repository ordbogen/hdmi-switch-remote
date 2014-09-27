package main

// http://www.atlona.com/pdf/manuals/AT-PRO3HD44M.pdf
// PWON PWON Power on
// PWOFF PWOFF Power off
// PWSTA PWx Will display the power status of the matrix (ex. Power is on = PWON)
// Version (Firmware #) Brings up the current firmware version
// Lock Lock Locks the front panel of the matrix so no buttons are active
// Unlock Unlock Unlocks the front panel of the matrix, enabling the buttons again
// All# x1AVx1,x2AVx2,... Resets all inputs to corresponding outputs (in3 to out3)
// x1$ x1$ Turns off output channel (to turn off output 3 = x3$)
// x1AVx2 x1AVx2 Switch input to output (input 3 to output 2 = x3AVx2)
// x1AVx2,x3,x4 x1AVx2,x3,x4 Switch input to multiple outputs (input 3 = x3AVx1,x2)
// RS232zoneX[command]
// ex: RS232zone1[PWON]
// RS232zoneX[command]
// ex: RS232zone1[PWON]
// Send commands to devices connected to the receiver RS-232 ports.
// Commands are the same as the ones stated in this table.
// X = zone number
// example: Turning the power on for the device connected in zone 1
// Statusx1 x7AVx1 Shows what input is connected to selected output
// Status x1AVx1,x2AVx2,
// x3AVx4, ....
// Displays which inputs are currently connected to which outputs
// SaveY (ex. Save2) SaveY (ex. Save2) Saves settings for future use, preset options 0 to 4
// RecallY (ex. Recall2) RecallY (ex. Recall2) Recalls saved settings for the number you selected
// ClearY (ex. Clear2) ClearY (ex. Clear2) Erases the save for the number you selected
// Mreset Mreset Sets matrix back to the default settings

const (
	CODE_POWER_ON     = "PWON"
	CODE_POWER_OFF    = "PWOFF"
	CODE_POWER_STATUS = "PWSTA"
	CODE_RESET        = "Mreset"
)
