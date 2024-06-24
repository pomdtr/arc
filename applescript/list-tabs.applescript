#!/usr/bin/osascript

 on escape_value(this_text)
  set AppleScript's text item delimiters to the "\\"
  set the item_list to every text item of this_text
  set AppleScript's text item delimiters to "\\\\"
  set this_text to the item_list as string
  set AppleScript's text item delimiters to the "\""
  set the item_list to every text item of this_text
  set AppleScript's text item delimiters to the "\\\""
  set this_text to the item_list as string
  set AppleScript's text item delimiters to ""
  return this_text
end escape_value

set _output to ""

tell application "Arc"
  tell first window
    set allTabs to properties of every tab
  end tell
  set tabsCount to count of allTabs
  repeat with i from 1 to tabsCount
    set _tab to item i of allTabs
    set _title to my escape_value(get title of _tab)
    set _url to get URL of _tab
    set _id to get id of _tab
    set _location to get location of _tab

    set _output to (_output & "{ \"title\": \"" & _title & "\", \"url\": \"" & _url & "\", \"id\": \"" & _id & "\", \"location\": \"" & _location & "\" }")

    if i < tabsCount then
      set _output to (_output & ",\n")
    else
      set _output to (_output & "\n")
    end if

  end repeat
end tell

return "[\n" & _output & "\n]"
