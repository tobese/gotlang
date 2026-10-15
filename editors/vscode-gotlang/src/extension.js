const vscode = require('vscode');
const path = require('path');
const fs = require('fs');

let gotlangTerminal;

function activate(context) {
  const disposable = vscode.commands.registerCommand('gotlang.runFile', runActiveFile);
  context.subscriptions.push(disposable);
}

function deactivate() {}

async function runActiveFile() {
  const editor = vscode.window.activeTextEditor;
  if (!editor) {
    vscode.window.showErrorMessage('Open a .got file first.');
    return;
  }

  const doc = editor.document;
  if (doc.languageId !== 'gotlang') {
    vscode.window.showErrorMessage('Active file is not a Gotlang (.got) file.');
    return;
  }

  if (doc.isDirty) {
    await doc.save();
  }

  const filePath = doc.fileName;
  const moduleRoot = findGotlangModuleRoot(path.dirname(filePath));

  const terminal = getOrCreateTerminal();
  terminal.show(true);

  if (moduleRoot) {
    // Developing inside (or against a checkout of) the gotlang repo itself:
    // run the CLI straight from source, no install required.
    const cliDir = path.join(moduleRoot, 'cmd', 'gotlang');
    terminal.sendText(`go run "${cliDir}" run "${filePath}"`);
  } else {
    // Otherwise assume a built "gotlang" binary is on PATH.
    terminal.sendText(`gotlang run "${filePath}"`);
  }
}

// findGotlangModuleRoot walks upward from startDir looking for the gotlang
// repo's cmd/gotlang/main.go, so "go run" works without any configuration
// when editing .got files inside (or under) a gotlang checkout.
function findGotlangModuleRoot(startDir) {
  let dir = startDir;
  for (;;) {
    if (fs.existsSync(path.join(dir, 'cmd', 'gotlang', 'main.go'))) {
      return dir;
    }
    const parent = path.dirname(dir);
    if (parent === dir) {
      return null;
    }
    dir = parent;
  }
}

function getOrCreateTerminal() {
  if (!gotlangTerminal || gotlangTerminal.exitStatus !== undefined) {
    gotlangTerminal = vscode.window.createTerminal('Gotlang');
  }
  return gotlangTerminal;
}

module.exports = { activate, deactivate };
