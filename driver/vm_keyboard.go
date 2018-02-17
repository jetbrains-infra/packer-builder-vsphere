package driver

import (
	"strings"
	"unicode/utf8"
	"unicode"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/methods"
)

func (vm *VirtualMachine) PutUsbScanCodes(s string) (int32, error) {
	spec := &types.PutUsbScanCodes{
		This: vm.vm.Reference(),
		Spec: getUsbScanCodeSpec(s),
	}
	resp, err := methods.PutUsbScanCodes(vm.driver.ctx, vm.driver.client.RoundTripper, spec)
	if err != nil {
		return 0, err
	}

	return resp.Returnval, nil
}

func getUsbScanCodeSpec(message string) types.UsbScanCodeSpec {
	scancodeIndex := make(map[string]uint)
	scancodeIndex["abcdefghijklmnopqrstuvwxyz"] = 4
	scancodeIndex["ABCDEFGHIJKLMNOPQRSTUVWXYZ"] = 4
	scancodeIndex["1234567890"] = 30
	scancodeIndex["!@#$%^&*()"] = 30
	scancodeIndex[" "] = 44
	scancodeIndex["-=[]\\"] = 45
	scancodeIndex["_+{}|" ] = 45
	scancodeIndex[ ";'`,./" ] = 51
	scancodeIndex[":\"~<>?" ] = 51

	shiftedChars := "!@#$%^&*()_+{}|:\"~<>?"

	scancodeMap := make(map[rune]uint)
	for chars, start := range scancodeIndex {
		var i uint = 0
		for len(chars) > 0 {
			r, size := utf8.DecodeRuneInString(chars)
			chars = chars[size:]
			scancodeMap[r] = start + i
			i += 1
		}
	}

	special := map[string]uint{
		"<enter>":    40,
		"<esc>":      41,
		"<bs>":       42,
		"<del>":      76,
		"<tab>":      43,
		"<f1>":       58,
		"<f2>":       59,
		"<f3>":       60,
		"<f4>":       61,
		"<f5>":       62,
		"<f6>":       63,
		"<f7>":       64,
		"<f8>":       65,
		"<f9>":       66,
		"<f10>":      67,
		"<f11>":      68,
		"<f12>":      69,
		"<insert>":   73,
		"<home>":     74,
		"<end>":      77,
		"<pageUp>":   75,
		"<pageDown>": 78,
		"<left>":     80,
		"<right>":    79,
		"<up>":       82,
		"<down>":     81,
	}

	spec := types.UsbScanCodeSpec{
		KeyEvents: []types.UsbScanCodeSpecKeyEvent{},
	}

	var keyAlt bool
	var keyCtrl bool
	var keyShift bool

	for len(message) > 0 {
		var scancode uint

		if strings.HasPrefix(message, "<leftAltOn>") {
			keyAlt = true
			message = message[len("<leftAltOn>"):]
		}

		if strings.HasPrefix(message, "<leftAltOff>") {
			keyAlt = false
			message = message[len("<leftAltOff>"):]
		}

		if strings.HasPrefix(message, "<leftCtrlOn>") {
			keyCtrl = true
			message = message[len("<leftCtrlOn>"):]
		}

		if strings.HasPrefix(message, "<leftCtrlOff>") {
			keyCtrl = false
			message = message[len("<leftCtrlOff>"):]
		}

		if strings.HasPrefix(message, "<leftShiftOn>") {
			keyShift = true
			message = message[len("<leftShiftOn>"):]
		}

		if strings.HasPrefix(message, "<leftShiftOff>") {
			keyShift = false
			message = message[len("<leftShiftOff>"):]
		}

		if scancode == 0 {
			for specialCode, specialValue := range special {
				if strings.HasPrefix(message, specialCode) {
					scancode = specialValue
					message = message[len(specialCode):]
					break
				}
			}
		}

		modAlt := keyAlt
		modCtrl := keyCtrl
		modShift := keyShift

		if scancode == 0 {
			r, size := utf8.DecodeRuneInString(message)
			scancode = scancodeMap[r]
			modShift = modShift || unicode.IsUpper(r) || strings.ContainsRune(shiftedChars, r)
			message = message[size:]
		}

		event := types.UsbScanCodeSpecKeyEvent{
			// https://github.com/lamw/vghetto-scripts/blob/f74bc8ba20064f46592bcce5a873b161a7fa3d72/powershell/VMKeystrokes.ps1#L130
			UsbHidCode: int32(scancode)<<16 | 7,
			Modifiers: &types.UsbScanCodeSpecModifierType{
				LeftAlt:     &modAlt,
				LeftControl: &modCtrl,
				LeftShift:   &modShift,
			},
		}
		spec.KeyEvents = append(spec.KeyEvents, event)
	}

	return spec
}
