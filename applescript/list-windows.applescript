#!/usr/bin/osascript
set _output to ""

on escape_value(this_text)
  set AppleScript's text item delimiters to the "\""
  set the item_list to every text item of this_text
  set AppleScript's text item delimiters to the "\\\""
  set this_text to the item_list as string
  set AppleScript's text item delimiters to ""
  return this_text
end replace_chars

tell application "Arc"
  set _window_index to 1

  repeat with _window in windows
    set _title to my escape_value(get name of _window)

    set _output to (_output & "{ \"title\": \"" & _title & "\", \"id\": " & _window_index & " }")

    if _window_index < (count windows) then
      set _output to (_output & ",\n")
    else
      set _output to (_output & "\n")
    end if

    set _window_index to _window_index + 1
  end repeat
end tell

return "[\n" & _output & "\n]"
