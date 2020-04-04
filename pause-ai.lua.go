package main

var pauseailua = []byte(`
local gui = require 'gui'

Pauser = defclass(Pauser, gui.Screen)
function Pauser:render(dc)
    self:renderParent()
    dc:string(' PAUSE ', COLOR_LIGHTCYAN)
end
function Pauser:onInput(keys)
    if keys.A_RETURN_TO_ARENA then
        self:dismiss()
    end
end

Pauser{}:show()
`)[1:]
