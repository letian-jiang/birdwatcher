package states

import (
	"strings"

	"github.com/spf13/cobra"
)

// State is the interface for application state.
type State interface {
	Label() string
	Process(cmd string) (State, error)
	Close()
	SetNext(state State)
}

// cmdState is the basic state to process input command.
type cmdState struct {
	label     string
	rootCmd   *cobra.Command
	nextState State
}

// Label returns the display label for current cli.
func (s *cmdState) Label() string {
	return s.label
}

// Process is the main entry for processing command.
func (s *cmdState) Process(cmd string) (State, error) {
	args := strings.Split(cmd, " ")

	s.rootCmd.SetArgs(args)
	err := s.rootCmd.Execute()
	if err != nil {
		return s, err
	}
	if s.nextState != nil {
		defer s.Close()
		// TODO fix ugly type cast
		if _, ok := s.nextState.(*exitState); ok {
			return s.nextState, ExitErr
		}
		return s.nextState, nil
	}
	// clean up args
	s.rootCmd.SetArgs(nil)
	return s, nil
}

// SetNext simple method to set next state.
func (s *cmdState) SetNext(state State) {
	s.nextState = state
}

// Close empty method to implement State.
func (s *cmdState) Close() {}

// Start returns the first state - offline.
func Start() State {
	root := &cobra.Command{
		Use:   "",
		Short: "",
	}

	state := &cmdState{
		label:   "Offline",
		rootCmd: root,
	}

	root.AddCommand(getConnectCommand(state), getExitCmd(state))
	return state
}
