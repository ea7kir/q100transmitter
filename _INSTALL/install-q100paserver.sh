#!/bin/bash

# Check current user
whoami | grep -q pi
if [ $? != 0 ]; then
    echo
    echo "Install must be performed as user pi"
    echo
    exit
fi

APPNAME="q100paserver"
Q100="/home/pi/AAA"

FILE=$Q100
if [ -e $Q100 ]; then
    echo
    echo "A Q100 app is already installed. If you wish to proceed, you will need"
    echo "to romove the SD Card and reinstall the OS using Raspberry Pi Imager."
    echo
    exit
fi

echo
echo "You are about to install $APPNAME into $HOME/Q100"
echo
echo "During the installation, the Raspberry Pi will restard, so"
echo "please BE PATIENT and wait until you see INSTALLATION COMPLETE"
echo

while true; do
    read -p "Do you wish to continue? (y/n) " yn
    case $yn in
        [Yy]* )
            echo
            echo "-----------------------------------"
            echo "----- Installing $APPNAME"
            echo "-----------------------------------"
            echo
            break;;
        [Nn]* )
            echo "Installation of $APPNAME is cancelled"
            exit;;
        * )
            echo "Please answer yes or no.";;
    esac
done

# exit

STATUS_LOG="/home/pi/install.log"

# Determinate whether the log file exists ? get the status : set status0
if [[ -f $STATUS_LOG ]]
then
    CURRENT_STATUS="$(cat "$STATUS_LOG")"
else
    CURRENT_STATUS="stage0"
    echo "$CURRENT_STATUS : $(date)"
    echo "$CURRENT_STATUS" > "$STATUS_LOG"
    # You could reboot at this point,
    # but probably you want to do action_1 first
fi

# Define your actions as functions

action_1()
{
    # do the 1st action

    mkdir $Q100
    cd $Q100
    touch done_1

    CURRENT_STATUS="stage1"
    echo "$CURRENT_STATUS : $(date)"
    echo "$CURRENT_STATUS" > "$STATUS_LOG"
    sudo reboot
    #exit # You could reboot at this point
}

action_2()
{
    # do the 2nd action
    cd $Q100
    touch done_2

    CURRENT_STATUS="stage2"
    echo "$CURRENT_STATUS : $(date)"
    echo "$CURRENT_STATUS" > "$STATUS_LOG"
    sudo reboot
    #exit # You could reboot at this point
}

action_3()
{
    # do the 3rd action
    cd $Q100
    touch done_3

    CURRENT_STATUS="stage3"
    echo "$CURRENT_STATUS : $(date)"
    echo "$CURRENT_STATUS" > "$STATUS_LOG"
    sudo reboot
    #exit # You could reboot at this point
}

case "$CURRENT_STATUS" in
stage0)
    action_1
    ;;
stage1)
    action_2
    ;;
stage2)
    action_3
    ;;
stage3)
    cd
    echo "---------------------------------"
    echo "----- INSTALLATION COMPLETE -----"
    echo "---------------------------------"
    echo
    ;;
*)
    echo "Something went wrong!"
    ;;
