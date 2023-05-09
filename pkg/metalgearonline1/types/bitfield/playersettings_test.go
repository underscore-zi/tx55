package bitfield

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestPlayerSettings_Byte1(t *testing.T) {
	var b PlayerSettings
	b.SetNameTags(true)
	if !b.GetNameTags() {
		t.Error("GetNameTags() should return true")
	}

	buf := bytes.Buffer{}
	_ = binary.Write(&buf, binary.BigEndian, b.Data)
	if buf.Bytes()[0] != 0x10 {
		t.Error("GetNameTags() should return true")
	}
}

func TestPlayerSettings_Byte2(t *testing.T) {
	var b PlayerSettings

	// First Nibble is just the switch speed
	b.SetSwitchSpeed(0x0F)
	if b.GetSwitchSpeed() != 0x0F {
		t.Error("GetSwitchSpeed() should return 0x0F")
	}

	// Second Nibble is the first person camera settings
	b.SetFirstPersonCameraSettings(true, true, true)
	if v, h, o := b.GetFirstPersonCameraSettings(); !v || !h || !o {
		t.Error("GetFirstPersonCameraSettings() should return true, true, true")
	}

	if b.Data[1].Data != 0xF7 {
		t.Error("Second byte should be 0xF7")
	}
}

func TestPlayerSettings_Byte3(t *testing.T) {
	var b PlayerSettings

	// First Nibble is the third person camera settings
	b.SetThirdPersonCameraSettings(true, true, true)
	if v, h, c := b.GetThirdPersonCameraSettings(); !v || !h || !c {
		t.Error("GetThirdPersonCameraSettings() should return true, true, true")
	}

	// Second Nibble is the first person camera rotation speed
	b.SetFirstPersonCameraRotationSpeed(0x0F)
	if b.GetFirstPersonCameraRotationSpeed() != 0x0F {
		t.Error("GetFirstPersonCameraRotationSpeed() should return 0x0F")
	}

	if b.Data[2].Data != 0x7F {
		t.Error("Third byte should be 0xF7")
	}
}

func TestPlayerSettings_Byte4(t *testing.T) {
	var b PlayerSettings

	b.SetCameraRotationSpeed(0x03)
	// First Nibble is the equipment switch style
	b.SetEquipmentSwitchStyle(0x03)
	if b.GetEquipmentSwitchStyle() != 0x03 {
		t.Error("SetEquipmentSwitchStyle() should return 0x0F")
	}

	if b.Data[3].Data != 0x33 {
		t.Error("Fourth byte should be 0x2F")
	}
}

func TestPlayerSettings_Byte5(t *testing.T) {
	var b PlayerSettings

	// First Nibble is the camera rotation speed
	b.SetUSBKeyboardType(0x02)
	if b.GetUSBKeyboardType() != 0x02 {
		t.Error("GetUSBKeyboardType() should return 0x02")
	}

	b.SetWeaponSwitchStyle(0x02)
	if b.GetWeaponSwitchStyle() != 0x02 {
		t.Error("GetWeaponSwitchStyle() should return 0x02")
	}

	b.SetUSBKeyboardType(0x02)
	if b.GetUSBKeyboardType() != 0x02 {
		t.Error("GetUSBKeyboardType() should return 0x02")
	}

	if b.Data[4].Data != 0x22 {
		t.Error("Fifth byte should be 0x0F")
	}
}
