package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	"runtime/trace"
	"strings"
)

func Profile(fn func() error) (err error) {
	if os.Getenv("TRACE") == "1" {
		err = startTraceProfiling()
		if err != nil {
			return err
		}
	}

	if os.Getenv("PROFILE_CPU") == "1" {
		err = startCPUProfiling()
		if err != nil {
			return err
		}
	}

	if err := fn(); err != nil {
		return err
	}

	for _, prof := range pprof.Profiles() {
		name := prof.Name()
		ev := "PROFILE_" + strings.ToUpper(name)

		if os.Getenv(ev) != "1" {
			continue
		}

		var fname string

		if v := os.Getenv(ev + "_FILE"); v != "" {
			fname = v
		} else {
			fname = name + ".pprof"
		}

		if err := writeProfileToFile(fname, name); err != nil {
			return fmt.Errorf("%w: %s", errCannotWriteProfile, name)
		}
	}

	return nil
}

var (
	errInvalidProfile        = errors.New("invalid profile given")
	errCannotCreateTraceFile = errors.New("cannot create trace execution file")
	errCannotStartTrace      = errors.New("cannot start execution tracing")
	errCannotCreateCPUFile   = errors.New("cannot create cpu profile file")
	errCannotProfileCPU      = errors.New("cannot profile cpu usage")
	errCannotWriteProfile    = errors.New("cannot write profile")
	errCannotCloseFile       = errors.New("cannot close file")
	errCannotCreateFile      = errors.New("cannot create file")
	errUnsafeFilename        = errors.New("unsafe filename")
)

func startTraceProfiling() (err error) {
	var fname string

	if v := os.Getenv("TRACE_FILE"); v != "" {
		fname = v
	} else {
		fname = "trace.out"
	}

	if !isSafeFilename(fname) {
		return fmt.Errorf("%w: %w", errUnsafeFilename, errors.New(fname))
	}

	//nolint:gosec // G304: The filename `fname` is checked for safety.
	f, err := os.Create(filepath.Join(".", fname)) //nolint:gosec
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotCreateTraceFile, err)
	}

	defer func() {
		if errC := f.Close(); errC != nil {
			errC = fmt.Errorf("%w: %w", errCannotCloseFile, errC)
			err = errors.Join(err, errC)
		}
	}()

	err = trace.Start(f)
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotStartTrace, err)
	}

	defer trace.Stop()

	return nil
}

func startCPUProfiling() (err error) {
	var fname string

	if v := os.Getenv("PROFILE_CPU_FILE"); v != "" {
		fname = v
	} else {
		fname = "cpu.pprof"
	}

	if !isSafeFilename(fname) {
		return fmt.Errorf("%w: %w", errUnsafeFilename, errors.New(fname))
	}

	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotCreateCPUFile, err)
	}

	defer func() {
		if errC := f.Close(); errC != nil {
			errC = fmt.Errorf("%w: %w", errCannotCloseFile, errC)
			err = errors.Join(err, errC)
		}
	}()

	err = pprof.StartCPUProfile(f)
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotProfileCPU, err)
	}

	defer pprof.StopCPUProfile()

	return nil
}

func writeProfile(w io.Writer, name string) error {
	prof := pprof.Lookup(name)
	if prof == nil {
		return errInvalidProfile
	}

	if err := prof.WriteTo(w, 0); err != nil {
		return fmt.Errorf("%w: %w", errCannotWriteProfile, err)
	}

	return nil
}

func writeProfileToFile(fname, name string) error {
	if !isSafeFilename(fname) {
		return fmt.Errorf("%w: %w", errUnsafeFilename, errors.New(fname))
	}

	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotCreateFile, err)
	}

	if err := writeProfile(f, name); err != nil {
		_ = f.Close() // Close the file even if writeProfile fails
		return err
	}

	return fmt.Errorf("%w: %w", errCannotCloseFile, f.Close())
}

func isSafeFilename(filename string) bool {
	cleaned := filepath.Clean(filename)
	return !strings.ContainsAny(cleaned, string(os.PathSeparator)) && !strings.HasPrefix(cleaned, "..") && !filepath.IsAbs(cleaned) && cleaned != "" && cleaned != "."
}
