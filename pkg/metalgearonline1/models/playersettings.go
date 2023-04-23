package models

import (
	"fmt"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/metalgearonline1/types/bitfield"
)

func init() {
	All = append(All, &PlayerSettings{})
}

type PlayerSettings struct {
	gorm.Model
	UserID uint

	ShowNameTags         bool
	Unknown              bool
	SwitchSpeed          byte
	FPVVertical          bool
	FPVHorizontal        bool
	FPVSwitchOrientation bool
	TPVVertical          bool
	TPVHorizontal        bool
	TPVChase             bool
	FPVRotationSpeed     byte
	EquipmentSwitchStyle byte
	TPVRotationSpeed     byte
	USBKeyboardType      byte
	WeaponSwitchStyle    byte
	FKey0                []byte
	FKey1                []byte
	FKey2                []byte
	FKey3                []byte
	FKey4                []byte
	FKey5                []byte
	FKey6                []byte
	FKey7                []byte
	FKey8                []byte
	FKey9                []byte
	FKey10               []byte
	FKey11               []byte
}

func (p *PlayerSettings) FromBitfield(b bitfield.PlayerSettings) {
	p.ShowNameTags = b.GetNameTags()
	p.Unknown = b.GetUnknown()
	p.SwitchSpeed = b.GetSwitchSpeed()
	p.FPVVertical, p.FPVHorizontal, p.FPVSwitchOrientation = b.GetFirstPersonCameraSettings()
	p.TPVVertical, p.TPVHorizontal, p.TPVChase = b.GetThirdPersonCameraSettings()
	p.FPVRotationSpeed = b.GetFirstPersonCameraRotationSpeed()
	p.TPVRotationSpeed = b.GetCameraRotationSpeed()
	p.EquipmentSwitchStyle = b.GetEquipmentSwitchStyle()
	p.USBKeyboardType = b.GetUSBKeyboardType()
	p.WeaponSwitchStyle = b.GetWeaponSwitchStyle()

	newbf := p.Bitfield()
	for i, cur := range newbf.Data {
		if cur.Data != b.Data[i].Data {
			// We could return an error here instead of printing, but really this "error" isn't something that should
			// ever happen. Its not something one can "recover" or do something about. Its just important if it happens
			// to know about it and fix it.
			fmt.Printf("bitfield mismatch [%d]: %x != %x\n", i, cur.Data, b.Data[i].Data)
		}
	}

}

func (p *PlayerSettings) Bitfield() bitfield.PlayerSettings {
	var b bitfield.PlayerSettings
	b.SetNameTags(p.ShowNameTags)
	b.SetUnknown(p.Unknown)
	b.SetSwitchSpeed(p.SwitchSpeed)
	b.SetFirstPersonCameraSettings(p.FPVVertical, p.FPVHorizontal, p.FPVSwitchOrientation)
	b.SetThirdPersonCameraSettings(p.TPVVertical, p.TPVHorizontal, p.TPVChase)
	b.SetFirstPersonCameraRotationSpeed(p.FPVRotationSpeed)
	b.SetCameraRotationSpeed(p.TPVRotationSpeed)
	b.SetEquipmentSwitchStyle(p.EquipmentSwitchStyle)
	b.SetUSBKeyboardType(p.USBKeyboardType)
	b.SetWeaponSwitchStyle(p.WeaponSwitchStyle)
	return b
}

func (p *PlayerSettings) PlayerSpecificSettings() types.PlayerSpecificSettings {
	var out types.PlayerSpecificSettings

	for i, cur := range p.Bitfield().Data {
		out.Settings.Data[i].Data = cur.Data
	}

	copy(out.FKeys[0][:], p.FKey0)
	copy(out.FKeys[1][:], p.FKey1)
	copy(out.FKeys[2][:], p.FKey2)
	copy(out.FKeys[3][:], p.FKey3)
	copy(out.FKeys[4][:], p.FKey4)
	copy(out.FKeys[5][:], p.FKey5)
	copy(out.FKeys[6][:], p.FKey6)
	copy(out.FKeys[7][:], p.FKey7)
	copy(out.FKeys[8][:], p.FKey8)
	copy(out.FKeys[9][:], p.FKey9)
	copy(out.FKeys[10][:], p.FKey10)
	copy(out.FKeys[11][:], p.FKey11)

	return out
}
