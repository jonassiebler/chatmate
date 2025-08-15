#!/usr/bin/env node

/**
 * ChatMate CLI - Transform VS Code Copilot Chat with specialized AI agents
 * 
 * This CLI tool allows global installation and management of chatmate agents
 * for VS Code Copilot Chat, providing a professional command-line interface
 * with auto-update capabilities.
 */

const { program } = require('commander');
const chalk = require('chalk');
const pkg = require('../package.json');
const { ChatMateManager } = require('../lib/chatmate-manager');

// Import update-notifier and handle gracefully
let notifier;

try {
  const updateNotifier = require('update-notifier');
  notifier = updateNotifier({ 
    pkg,
    updateCheckInterval: 1000 * 60 * 60 * 24 // Check once per day
  });
} catch (error) {
  // Gracefully handle update-notifier issues
  notifier = null;
}

// Notify about updates
if (notifier && notifier.update) {
  console.log(chalk.yellow(`
╭─────────────────────────────────────────────────────────╮
│                                                         │
│    Update available: ${chalk.green(notifier.update.latest)} (current: ${chalk.red(notifier.update.current)})    │
│                                                         │
│    Run ${chalk.cyan('npm install -g chatmate')} to update.               │
│                                                         │
╰─────────────────────────────────────────────────────────╯
  `));
}

// Initialize ChatMate manager
const chatmate = new ChatMateManager();

// Configure CLI program
program
  .name('chatmate')
  .description('Open source collection of specialized AI agents for VS Code Copilot Chat')
  .version(pkg.version, '-v, --version', 'display version number')
  .helpOption('-h, --help', 'display help for command');

// Main hire command - preserves hire.sh functionality
program
  .command('hire')
  .description('Install all chatmate agents to VS Code Copilot Chat')
  .option('-s, --specific <agents...>', 'install only specific chatmate agents')
  .option('-f, --force', 'force reinstall all agents (overwrite existing)')
  .action(async (options) => {
    try {
      console.log(chalk.blue('🤖 Installing chatmate agents...\n'));
      
      if (options.specific) {
        await chatmate.installSpecific(options.specific, options.force);
      } else {
        await chatmate.installAll(options.force);
      }
      
      console.log(chalk.green('\n✅ All chatmates installed! Restart VS Code to use them.'));
    } catch (error) {
      console.error(chalk.red(`❌ Error: ${error.message}`));
      process.exit(1);
    }
  });

// List command - show available and installed chatmates
program
  .command('list')
  .alias('ls')
  .description('List available and installed chatmate agents')
  .option('-a, --available', 'show only available chatmates')
  .option('-i, --installed', 'show only installed chatmates')
  .action(async (options) => {
    try {
      await chatmate.listChatmates(options);
    } catch (error) {
      console.error(chalk.red(`❌ Error: ${error.message}`));
      process.exit(1);
    }
  });

// Uninstall command - remove specific chatmates
program
  .command('uninstall')
  .alias('remove')
  .description('Uninstall specific chatmate agents')
  .argument('<agents...>', 'chatmate agents to uninstall')
  .option('-a, --all', 'uninstall all chatmate agents')
  .action(async (agents, options) => {
    try {
      console.log(chalk.yellow('🗑️  Uninstalling chatmate agents...\n'));
      
      if (options.all) {
        await chatmate.uninstallAll();
      } else {
        await chatmate.uninstallSpecific(agents);
      }
      
      console.log(chalk.green('\n✅ Chatmates uninstalled successfully!'));
    } catch (error) {
      console.error(chalk.red(`❌ Error: ${error.message}`));
      process.exit(1);
    }
  });

// Status command - show VS Code and chatmate installation status
program
  .command('status')
  .description('Show ChatMate and VS Code installation status')
  .action(async () => {
    try {
      await chatmate.showStatus();
    } catch (error) {
      console.error(chalk.red(`❌ Error: ${error.message}`));
      process.exit(1);
    }
  });

// Config command - manage ChatMate configuration
program
  .command('config')
  .description('Manage ChatMate configuration settings')
  .option('--show', 'show current configuration')
  .option('--reset', 'reset configuration to defaults')
  .action(async (options) => {
    try {
      await chatmate.manageConfig(options);
    } catch (error) {
      console.error(chalk.red(`❌ Error: ${error.message}`));
      process.exit(1);
    }
  });

// Parse command line arguments
program.parse();

// Show help if no command provided
if (!process.argv.slice(2).length) {
  program.outputHelp();
}
