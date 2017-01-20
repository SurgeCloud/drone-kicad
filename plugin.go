package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	pythonexec       = "python2"
	sch_script       = "/bin/ci-scripts/export_schematic.py"
	bom_script       = "/bin/ci-scripts/export_bom.py"
	grb_script       = "/bin/ci-scripts/export_grb.py"
	fplibtablescript = "gen_fp_lib_table.sh"

	// Binary masks - MSBits
	badh_mask     = 1 << 0
	fadh_mask     = 1 << 1
	bpaste_mask   = 1 << 2
	fpaste_mask   = 1 << 3
	bsilks_mask   = 1 << 4
	fsilks_mask   = 1 << 5
	bmask_mask    = 1 << 6
	fmask_mask    = 1 << 7
	dwguser_mask  = 1 << 8
	cmtsuser_mask = 1 << 9
	eco1_mask     = 1 << 10
	eco2_mask     = 1 << 11
	ecuts_mask    = 1 << 12
	margin_mask   = 1 << 13
	bcrtyd_mask   = 1 << 14
	fcrtyd_mask   = 1 << 15
	bfab_mask     = 1 << 16
	ffab_mask     = 1 << 17

	// Binary masks - LSBits
	fcu_mask = 1 << 3
	bcu_mask = 1 << 31
)

type (
	// Client defines the client data to be embedded in some documents
	Client struct {
		Code string // Enterprise client code
		Name string // Enterprise client name
	}

	// Project defines the KiCad project
	Project struct {
		Code string // Enterprise project code
		Name string // Enterprise project name
	}

	// Gerber defines the options for exporting Gerber files
	GerberLayers struct {
		Fcu bool `json:"fcu"`
		Bcu bool `json:"bcu"`
		Fm  bool `json:"fm"`
		Bm  bool `json:"bm"`
		Fs  bool `json:"fs"`
		Bs  bool `json:"bs"`
		Ec  bool `json:"ec"`
	}

	// Options defines what to generate
	Options struct {
		Sch bool // Generate Schematic (pdf)
		Bom bool // Generate BOM (xml & xlsx)
		//Brd	bool // Generate PCB plot (pdf)
		Grb    GerberLayers // Gerber file layers
		GrbGen bool         // Generate Gerber files
		//Lyr	bool // Generate plot for each layer (pdf)
		//Wrl	bool // Generate VRML PCB
		//Stp	bool // Generate Step PCB
		//3d	bool // Generate plot of 3D view (png)
	}

	// Plugin defines the KiCad plugin parameters
	Plugin struct {
		Client  Client  // Client configuration
		Project Project // Project configuration
		Options Options // Plugin options
	}
)

func (p Plugin) Exec() error {

	var cmds []*exec.Cmd

	if p.Options.Sch {
		cmds = append(cmds, commandSchematic(p.Project))
	}
	if p.Options.Bom {
		cmds = append(cmds, commandBOM(p.Project))
	}
	if p.Options.GrbGen {
		cmds = append(cmds, commandGenFpLibTable())
		cmds = append(cmds, commandSetGerberLayers(p.Project, p.Options.Grb))
		cmds = append(cmds, commandGerber(p.Project))
	}

	// Set env variables
	os.Setenv("DISPLAY", ":0")
	os.Setenv("DEBIAN_FRONTEND", "noninteractive")

	// execute all commands in batch mode.
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func commandGenFpLibTable() *exec.Cmd {

	return exec.Command(
		fplibtablescript,
	)
}

func commandSetGerberLayers(pjt Project, lyr GerberLayers) *exec.Cmd {

	var layerselection_str string
	var lyr_lsb uint32 = 0
	var lyr_msb uint32 = 0
	if lyr.Fcu {
		lyr_lsb |= fcu_mask
	}
	if lyr.Bcu {
		lyr_lsb |= bcu_mask
	}
	if lyr.Fs {
		lyr_msb |= fsilks_mask
	}
	if lyr.Bs {
		lyr_msb |= bsilks_mask
	}
	if lyr.Fm {
		lyr_msb |= fmask_mask
	}
	if lyr.Bm {
		lyr_msb |= bmask_mask
	}

	layerselection_str = fmt.Sprintf("%#08x_%08x", lyr_msb, lyr_lsb)

	var sed_cmd string
	sed_cmd = fmt.Sprintf("%s %s%s", "s/\\([\\s\\t]*layerselection\\).*$/\\1", layerselection_str, ")/")

	var brd_name string
	brd_name = fmt.Sprintf("%s.kicad_pcb", pjt.Name)

	return exec.Command(
		"sed",
		"-i",
		sed_cmd,
		brd_name,
	)
}

func commandGerber(pjt Project) *exec.Cmd {

	return exec.Command(
		pythonexec,
		"-u",
		grb_script,
		pjt.Name,
	)
}

func commandSchematic(pjt Project) *exec.Cmd {

	return exec.Command(
		pythonexec,
		"-u",
		sch_script,
		pjt.Name,
	)
}

func commandBOM(pjt Project) *exec.Cmd {

	return exec.Command(
		pythonexec,
		"-u",
		bom_script,
		pjt.Name,
	)
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
