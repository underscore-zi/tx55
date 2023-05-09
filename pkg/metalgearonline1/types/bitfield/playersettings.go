package bitfield

import "fmt"

type PlayerSettings struct {
	Data [8]Bitfield
}

const (
	// Byte 1

	NameTagsPosition = 4
	UnknownPosition  = 0

	// Byte 2

	SwitchSpeedPosition           = 4
	SwitchSpeedSize               = 4
	FirstPersonVerticalPosition   = 0
	FirstPersonHorizontalPosition = 1
	FirstPersonSwitchOrientation  = 2

	// Byte 3

	ThirdPersonVerticalPosition      = 4
	ThirdPersonHorizontalPosition    = 5
	ChaseCameraWhileShootingPosition = 6
	FirstPersonRotationSpeedPosition = 0
	FirstPersonRotationSpeedSize     = 4

	// Byte 4

	EquipmentSwitchStylePosition = 4
	EquipmentSwitchStyleSize     = 2
	CameraRotationSpeedPosition  = 0
	CameraRotationSpeedSize      = 4

	// Byte 5

	USBKeyboardTypePosition   = 4
	USBKeyboardTypeSize       = 2
	WeaponSwitchStylePosition = 0
	WeaponSwitchStyleSize     = 2
)

// Byte 1

func (p *PlayerSettings) SetUnknown(value bool) {
	p.Data[0].SetBit(UnknownPosition, value)
}

func (p *PlayerSettings) GetUnknown() bool {
	return p.Data[0].GetBit(UnknownPosition)
}

func (p *PlayerSettings) SetNameTags(value bool) {
	p.Data[0].SetBit(NameTagsPosition, value)
}

func (p *PlayerSettings) GetNameTags() bool {
	return p.Data[0].GetBit(NameTagsPosition)
}

// Byte 2

func (p *PlayerSettings) SetSwitchSpeed(value byte) {
	p.Data[1].SetBits(SwitchSpeedPosition, SwitchSpeedSize, value)
}

func (p *PlayerSettings) GetSwitchSpeed() byte {
	return p.Data[1].GetBits(SwitchSpeedPosition, SwitchSpeedSize)
}

func (p *PlayerSettings) SetFirstPersonCameraSettings(vertical, horizontal, switchOrientation bool) {
	p.Data[1].SetBit(FirstPersonVerticalPosition, vertical)
	p.Data[1].SetBit(FirstPersonHorizontalPosition, horizontal)
	p.Data[1].SetBit(FirstPersonSwitchOrientation, switchOrientation)
}

func (p *PlayerSettings) GetFirstPersonCameraSettings() (vertical, horizontal, switchOrientation bool) {
	vertical = p.Data[1].GetBit(FirstPersonVerticalPosition)
	horizontal = p.Data[1].GetBit(FirstPersonHorizontalPosition)
	switchOrientation = p.Data[1].GetBit(FirstPersonSwitchOrientation)
	return
}

// Byte 3

func (p *PlayerSettings) SetThirdPersonCameraSettings(vertical, horizontal, chaseCamera bool) {
	p.Data[2].SetBit(ThirdPersonVerticalPosition, vertical)
	p.Data[2].SetBit(ThirdPersonHorizontalPosition, horizontal)
	p.Data[2].SetBit(ChaseCameraWhileShootingPosition, chaseCamera)
}

func (p *PlayerSettings) GetThirdPersonCameraSettings() (vertical, horizontal, chaseCamera bool) {
	vertical = p.Data[2].GetBit(ThirdPersonVerticalPosition)
	horizontal = p.Data[2].GetBit(ThirdPersonHorizontalPosition)
	chaseCamera = p.Data[2].GetBit(ChaseCameraWhileShootingPosition)
	return
}

func (p *PlayerSettings) SetFirstPersonCameraRotationSpeed(value byte) {
	p.Data[2].SetBits(FirstPersonRotationSpeedPosition, FirstPersonRotationSpeedSize, value)
}

func (p *PlayerSettings) GetFirstPersonCameraRotationSpeed() byte {
	return p.Data[2].GetBits(FirstPersonRotationSpeedPosition, FirstPersonRotationSpeedSize)
}

// Byte 4

func (p *PlayerSettings) SetEquipmentSwitchStyle(value byte) {
	p.Data[3].SetBits(EquipmentSwitchStylePosition, EquipmentSwitchStyleSize, value)
}

func (p *PlayerSettings) GetEquipmentSwitchStyle() byte {
	return p.Data[3].GetBits(EquipmentSwitchStylePosition, EquipmentSwitchStyleSize)
}

func (p *PlayerSettings) SetCameraRotationSpeed(value byte) {
	p.Data[3].SetBits(CameraRotationSpeedPosition, CameraRotationSpeedSize, value)
}

func (p *PlayerSettings) GetCameraRotationSpeed() byte {
	return p.Data[3].GetBits(CameraRotationSpeedPosition, CameraRotationSpeedSize)
}

// Byte 5

func (p *PlayerSettings) SetUSBKeyboardType(value byte) {
	p.Data[4].SetBits(USBKeyboardTypePosition, USBKeyboardTypeSize, value)
}

func (p *PlayerSettings) GetUSBKeyboardType() byte {
	return p.Data[4].GetBits(USBKeyboardTypePosition, USBKeyboardTypeSize)
}

func (p *PlayerSettings) SetWeaponSwitchStyle(value byte) {
	p.Data[4].SetBits(WeaponSwitchStylePosition, WeaponSwitchStyleSize, value)
}

func (p *PlayerSettings) GetWeaponSwitchStyle() byte {
	return p.Data[4].GetBits(WeaponSwitchStylePosition, WeaponSwitchStyleSize)
}

func (p *PlayerSettings) String() string {
	nameTags := "Off"
	if p.GetNameTags() {
		nameTags = "On"
	}

	firstPersonVertical, firstPersonHorizontal, firstPersonSwitchOrientation := p.GetFirstPersonCameraSettings()
	firstPersonVerticalStr := "Normal"
	if firstPersonVertical {
		firstPersonVerticalStr = "Reverse"
	}
	firstPersonHorizontalStr := "Normal"
	if firstPersonHorizontal {
		firstPersonHorizontalStr = "Reverse"
	}
	firstPersonSwitchOrientationStr := "Player Orientation"
	if firstPersonSwitchOrientation {
		firstPersonSwitchOrientationStr = "Camera Orientation"
	}

	thirdPersonVertical, thirdPersonHorizontal, chaseCamera := p.GetThirdPersonCameraSettings()
	thirdPersonVerticalStr := "Normal"
	if thirdPersonVertical {
		thirdPersonVerticalStr = "Reverse"
	}
	thirdPersonHorizontalStr := "Normal"
	if thirdPersonHorizontal {
		thirdPersonHorizontalStr = "Reverse"
	}
	chaseCameraStr := "Off"
	if chaseCamera {
		chaseCameraStr = "On"
	}

	return fmt.Sprintf("Name Tags: %s\nSwitch Speed: %d\nFirst Person Camera Settings: Vertical=%s, Horizontal=%s, Switch Orientation=%s\nThird Person Camera Settings: Vertical=%s, Horizontal=%s, Chase Camera=%s\nFirst Person Camera Rotation Speed: %d\nEquipment Switch Style: %d\nCamera Rotation Speed: %d\nUSB Keyboard Type: %d\nWeapon Switch Style: %d\n%x",
		nameTags,
		p.GetSwitchSpeed(),
		firstPersonVerticalStr,
		firstPersonHorizontalStr,
		firstPersonSwitchOrientationStr,
		thirdPersonVerticalStr,
		thirdPersonHorizontalStr,
		chaseCameraStr,
		p.GetFirstPersonCameraRotationSpeed(),
		p.GetEquipmentSwitchStyle(),
		p.GetCameraRotationSpeed(),
		p.GetUSBKeyboardType(),
		p.GetWeaponSwitchStyle(),
		p.Data)
}
