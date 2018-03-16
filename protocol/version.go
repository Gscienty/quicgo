package protocol

type Version		uint32

func (this Version) SupportedVersion(other Version) bool {
	return this == other
}

func (this Version) SupportedVersions(others []Version) bool {
	for _, v := range others {
		if v == this {
			return true
		}
	}
	return false
}

func ChooseSupportedVersion(ourVersions []Version, supportedVersion []Version) (Version, bool) {
	for _, ourVersion := range ourVersions {
		for _, theirVersion := range supportedVersion {
			if ourVersion == theirVersion {
				return ourVersion, true
			}
		}
	}

	return 0, false
}