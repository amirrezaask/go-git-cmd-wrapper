/*
* CODE GENERATED AUTOMATICALLY
* THIS FILE MUST NOT BE EDITED BY HAND
 */
package add

import "github.com/ldez/go-git-cmd-wrapper/types"

// All Update the index not only where the working tree has a file matching <pathspec> but also where the index already has an entry. 
// This adds, modifies, and removes index entries to match the working tree. 
// If no <pathspec> is given when -A option is used, all files in the entire working tree are updated (old versions of Git used to limit the update to the current directory and its subdirectories). 
// -A, --all, --no-ignore-removal
func All(g *types.Cmd) {
	g.AddOptions("--all")
}

// DryRun Don't actually add the file(s), just show if they exist and/or will be ignored. 
// -n, --dry-run
func DryRun(g *types.Cmd) {
	g.AddOptions("--dry-run")
}

// Edit Open the diff vs. 
// the index in an editor and let the user edit it. 
// After the editor was closed, adjust the hunk headers and apply the patch to the index. 
// The intent of this option is to pick and choose lines of the patch to apply, or even to modify the contents of lines to be staged. 
// This can be quicker and more flexible than using the interactive hunk selector. 
// However, it is easy to confuse oneself and create a patch that does not apply to the index. 
// See EDITING PATCHES below. 
// -e, --edit
func Edit(g *types.Cmd) {
	g.AddOptions("--edit")
}

// Force Allow adding otherwise ignored files. 
// -f, --force
func Force(g *types.Cmd) {
	g.AddOptions("--force")
}

// IgnoreErrors If some files could not be added because of errors indexing them, do not abort the operation, but continue adding the others. 
// The command shall still exit with non-zero status. 
// The configuration variable add.ignoreErrors can be set to true to make this the default behaviour. 
// --ignore-errors
func IgnoreErrors(g *types.Cmd) {
	g.AddOptions("--ignore-errors")
}

// IgnoreMissing This option can only be used together with --dry-run. 
// By using this option the user can check if any of the given files would be ignored, no matter if they are already present in the work tree or not. 
// --ignore-missing
func IgnoreMissing(g *types.Cmd) {
	g.AddOptions("--ignore-missing")
}

// IntentToAdd Record only the fact that the path will be added later. 
// An entry for the path is placed in the index with no content. 
// This is useful for, among other things, showing the unstaged content of such files with git diff and committing them with git commit -a. 
// -N, --intent-to-add
func IntentToAdd(g *types.Cmd) {
	g.AddOptions("--intent-to-add")
}

// Interactive Add modified contents in the working tree interactively to the index. 
// Optional path arguments may be supplied to limit operation to a subset of the working tree. 
// See 'Interactive mode' for details. 
// -i, --interactive
func Interactive(g *types.Cmd) {
	g.AddOptions("--interactive")
}

// NoAll Update the index by adding new files that are unknown to the index and files modified in the working tree, but ignore files that have been removed from the working tree. 
// This option is a no-op when no <pathspec> is used. 
// This option is primarily to help users who are used to older versions of Git, whose 'git add <pathspec>...' was a synonym for 'git add --no-all <pathspec>...', i.e. ignored removed files. 
// --no-all, --ignore-removal
func NoAll(g *types.Cmd) {
	g.AddOptions("--no-all")
}

// Patch Interactively choose hunks of patch between the index and the work tree and add them to the index. 
// This gives the user a chance to review the difference before adding modified contents to the index. 
// This effectively runs add --interactive, but bypasses the initial command menu and directly jumps to the patch subcommand. 
// See “Interactive mode” for details. 
// -p, --patch
func Patch(g *types.Cmd) {
	g.AddOptions("--patch")
}

// Refresh Don't add the file(s), but only refresh their stat() information in the index. 
// --refresh
func Refresh(g *types.Cmd) {
	g.AddOptions("--refresh")
}

// Update Update the index just where it already has an entry matching <pathspec>. 
// This removes as well as modifies index entries to match the working tree, but adds no new files. 
// If no <pathspec> is given when -u option is used, all tracked files in the entire working tree are updated (old versions of Git used to limit the update to the current directory and its subdirectories). 
// -u, --update
func Update(g *types.Cmd) {
	g.AddOptions("--update")
}

// Verbose Be verbose. 
// -v, --verbose
func Verbose(g *types.Cmd) {
	g.AddOptions("--verbose")
}
