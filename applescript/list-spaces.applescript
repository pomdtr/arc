#!/usr/bin/osascript

set _output to ""

tell application "Arc"
  set _space_index to 1

  repeat with _space in spaces of front window
    set _title to get title of _space

    set _output to (_output & "{ \"title\": \"" & _title & "\", \"id\": " & _space_index & " }")

    if _space_index < (count spaces of front window) then
      set _output to (_output & ",\n")
    else
      set _output to (_output & "\n")
    end if

    set _space_index to _space_index + 1
  end repeat
end tell

return "[\n" & _output & "\n]"

