// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {

	// Use the console to output diagnostic information (console.log) and errors (console.error)
	// This line of code will only be executed once when your extension is activated
	console.log('Congratulations, your extension "lzr" is now active!');

	// The command has been defined in the package.json file
	// Now provide the implementation of the command with registerCommand
	// The commandId parameter must match the command field in package.json
	const disposable = vscode.commands.registerCommand('lzr.helloWorld', () => {
		// The code you place here will be executed every time your command is executed
		// Display a message box to the user
		vscode.window.showInformationMessage('Hello World from lzr!');
	});

	// Register hover provider for Go-like variables in .lzr files
    const hoverProvider = vscode.languages.registerHoverProvider('lzr', {
		provideHover(document, position, token) {
			const wordRangeGoVar = document.getWordRangeAtPosition(position, /\{\{\s*[a-zA-Z_][a-zA-Z0-9_]*\.[a-zA-Z_][a-zA-Z0-9_]*\s*\}\}/);
			const wordRangeHtml = document.getWordRangeAtPosition(position);
			const wordHtml = document.getText(wordRangeHtml);
			
			// Check for Go variable hover
			if (wordRangeGoVar) {
				const goVariable = document.getText(wordRangeGoVar);
				return new vscode.Hover(`Go Variable: ${goVariable}`);
			}
	
			// Check for HTML element hover (this is basic; you can enhance it)
			if (wordHtml) {
				return new vscode.Hover(`HTML Element: ${wordHtml}`);
			}
	
			return undefined; // Return undefined if no hover is applicable
		}
	});
	

    // Register autocompletion provider for Go-like variables in .lzr files
    const completionProvider = vscode.languages.registerCompletionItemProvider('lzr', {
        provideCompletionItems(document, position, token, context) {
            const completionItem = new vscode.CompletionItem('Name', vscode.CompletionItemKind.Variable);
            completionItem.detail = 'A Go variable';
            completionItem.insertText = '{{ .Name }}';
            return [completionItem];
        }
    });

	context.subscriptions.push(hoverProvider);
    context.subscriptions.push(completionProvider);
	context.subscriptions.push(disposable);
}

// This method is called when your extension is deactivated
export function deactivate() {}
