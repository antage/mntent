package mntent

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func setup(lines []string, t *testing.T) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "mntent_test")
	if err != nil {
		return "", err
	}
	defer file.Close()

	tempFileName := file.Name()

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	err = file.Truncate(0)
	if err != nil {
		return "", err
	}

	for _, line := range lines {
		_, err := file.WriteString(line)
		if err != nil {
			return "", err
		}
		_, err = file.WriteString("\n")
		if err != nil {
			return "", err
		}
	}

	return tempFileName, nil
}

func teardown(filename string) {
	os.Remove(filename)
}

func TestFstab(t *testing.T) {
	data := []string{
		"# /etc/fstab: static file system information.",
		"#",
		"# Use 'blkid' to print the universally unique identifier for a",
		"# device; this may be used with UUID= as a more robust way to name devices",
		"# that works even if disks are added and removed. See fstab(5).",
		"#",
		"# <file system>                         <mount point>   <type>  <options>                           <dump>  <pass>",
		"/dev/mapper/vgssd-root                  /               ext4    noatime,discard,errors=remount-ro   0       1",
		"/dev/mapper/vg1-data                    /data           xfs     noatime,nodev,nosuid,noexec         0       2",
		"/dev/mapper/vg--raid1-home              /home           xfs     defaults                            0       2",
		"/dev/mapper/vg--raid1-var               /var            ext4    defaults                            0       2",
		"",
		"/dev/mapper/swap1			none            swap    sw                                  0       0",
		"/dev/mapper/swap2			none            swap    sw                                  0       0",
		"",
		"/dev/sr0                                /media/cdrom0   udf,iso9660 user,noauto                     0       0",
		"",
		"cgroup                                  /sys/fs/cgroup  cgroup  defaults                            0       0",
	}
	tabFileName, err := setup(data, t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer teardown(tabFileName)

	entries, err := Parse(tabFileName)
	if err != nil {
		t.Fatalf("Can't parse %s: %s", tabFileName, err)
	}

	expected := []Entry{
		Entry{
			Name:          "/dev/mapper/vgssd-root",
			Directory:     "/",
			Types:         []string{"ext4"},
			Options:       []string{"noatime", "discard", "errors=remount-ro"},
			DumpFrequency: 0,
			PassNumber:    1,
		},
		Entry{
			Name:          "/dev/mapper/vg1-data",
			Directory:     "/data",
			Types:         []string{"xfs"},
			Options:       []string{"noatime", "nodev", "nosuid", "noexec"},
			DumpFrequency: 0,
			PassNumber:    2,
		},
		Entry{
			Name:          "/dev/mapper/vg--raid1-home",
			Directory:     "/home",
			Types:         []string{"xfs"},
			Options:       []string{"defaults"},
			DumpFrequency: 0,
			PassNumber:    2,
		},
		Entry{
			Name:          "/dev/mapper/vg--raid1-var",
			Directory:     "/var",
			Types:         []string{"ext4"},
			Options:       []string{"defaults"},
			DumpFrequency: 0,
			PassNumber:    2,
		},
		Entry{
			Name:          "/dev/mapper/swap1",
			Directory:     "none",
			Types:         []string{"swap"},
			Options:       []string{"sw"},
			DumpFrequency: 0,
			PassNumber:    0,
		},
		Entry{
			Name:          "/dev/mapper/swap2",
			Directory:     "none",
			Types:         []string{"swap"},
			Options:       []string{"sw"},
			DumpFrequency: 0,
			PassNumber:    0,
		},
		Entry{
			Name:          "/dev/sr0",
			Directory:     "/media/cdrom0",
			Types:         []string{"udf", "iso9660"},
			Options:       []string{"user", "noauto"},
			DumpFrequency: 0,
			PassNumber:    0,
		},
		Entry{
			Name:          "cgroup",
			Directory:     "/sys/fs/cgroup",
			Types:         []string{"cgroup"},
			Options:       []string{"defaults"},
			DumpFrequency: 0,
			PassNumber:    0,
		},
	}

	if len(entries) != len(expected) {
		t.Fatalf("Expected %d filesystems but got %d", len(expected), len(entries))
	}

	for i, exp := range expected {
		t.Logf("Compare data for %q filesystem.", exp.Name)
		if entries[i].Name != exp.Name {
			t.Errorf("Expected %q filesystem but got %q", exp.Name, entries[i].Name)
		}
		if entries[i].Directory != exp.Directory {
			t.Errorf("Expected %q directory but got %q", exp.Directory, entries[i].Directory)
		}
		if !reflect.DeepEqual(entries[i].Types, exp.Types) {
			t.Errorf("Expected %v filesystem type but got %v", exp.Types, entries[i].Types)
		}
		if !reflect.DeepEqual(entries[i].Options, exp.Options) {
			t.Errorf("Expected %v filesystem options but got %v", exp.Options, entries[i].Options)
		}
		if entries[i].DumpFrequency != exp.DumpFrequency {
			t.Errorf("Expected %d dump frequency but got %d", exp.DumpFrequency, entries[i].DumpFrequency)
		}
		if entries[i].PassNumber != exp.PassNumber {
			t.Errorf("Expected %d pass number but got %d", exp.PassNumber, entries[i].PassNumber)
		}
	}

}

func ExampleFstab() {
	entries, err := Parse("/etc/fstab")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open /etc/fstab: %s\n", err)
		os.Exit(1)
	}

	for _, ent := range entries {
		fmt.Printf("FS name: %s\n", ent.Name)
		fmt.Printf("\tDir: %s\n", ent.Directory)
		fmt.Printf("\tTypes: %v\n", ent.Types)
		fmt.Printf("\tOptions: %s\n", strings.Join(ent.Options, ","))
		fmt.Printf("\tDump frequency: %d\n", ent.DumpFrequency)
		fmt.Printf("\tFSCK pass number: %d\n", ent.PassNumber)
	}
}
