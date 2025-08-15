const fs = require('fs-extra');
const path = require('path');
const os = require('os');
const chalk = require('chalk');

/**
 * ChatMateManager - Core logic for managing chatmate agents
 * 
 * This class provides the core functionality for installing, uninstalling,
 * and managing chatmate agents for VS Code Copilot Chat. It replaces the
 * functionality of the hire.sh script with enhanced Node.js capabilities.
 */
class ChatMateManager {
  constructor() {
    this.scriptDir = path.resolve(__dirname, '..');
    this.matesDir = path.join(this.scriptDir, 'mates');
    this.promptsDir = this.getPromptsDirectory();
  }

  /**
   * Get the VS Code Copilot Chat prompts directory based on OS
   * @returns {string} The prompts directory path
   */
  getPromptsDirectory() {
    const platform = os.platform();
    const homeDir = os.homedir();
    
    switch (platform) {
      case 'darwin': // macOS
        return path.join(homeDir, 'Library', 'Application Support', 'Code', 'User', 'prompts');
      
      case 'linux': // Linux
        return path.join(homeDir, '.config', 'Code', 'User', 'prompts');
      
      case 'win32': // Windows
        const appData = process.env.APPDATA || path.join(homeDir, 'AppData', 'Roaming');
        return path.join(appData, 'Code', 'User', 'prompts');
      
      default:
        throw new Error(`Unsupported operating system: ${platform}`);
    }
  }

  /**
   * Get all available chatmate files
   * @returns {Promise<string[]>} Array of chatmate filenames
   */
  async getAvailableChatmates() {
    try {
      const files = await fs.readdir(this.matesDir);
      return files.filter(file => file.endsWith('.chatmode.md'));
    } catch (error) {
      throw new Error(`Failed to read chatmates directory: ${error.message}`);
    }
  }

  /**
   * Get all installed chatmate files
   * @returns {Promise<string[]>} Array of installed chatmate filenames
   */
  async getInstalledChatmates() {
    try {
      if (!await fs.pathExists(this.promptsDir)) {
        return [];
      }
      
      const files = await fs.readdir(this.promptsDir);
      return files.filter(file => file.endsWith('.chatmode.md'));
    } catch (error) {
      throw new Error(`Failed to read prompts directory: ${error.message}`);
    }
  }

  /**
   * Install all chatmate agents
   * @param {boolean} force - Force reinstall even if already installed
   */
  async installAll(force = false) {
    try {
      // Ensure prompts directory exists
      await fs.ensureDir(this.promptsDir);
      
      const chatmates = await this.getAvailableChatmates();
      
      if (chatmates.length === 0) {
        throw new Error('No chatmate files found in mates directory');
      }

      console.log(`Installing ${chatmates.length} chatmates to: ${chalk.cyan(this.promptsDir)}\n`);

      for (const chatmate of chatmates) {
        await this.installChatmate(chatmate, force);
      }

    } catch (error) {
      throw new Error(`Failed to install chatmates: ${error.message}`);
    }
  }

  /**
   * Install specific chatmate agents
   * @param {string[]} agentNames - Array of agent names to install
   * @param {boolean} force - Force reinstall even if already installed
   */
  async installSpecific(agentNames, force = false) {
    try {
      // Ensure prompts directory exists
      await fs.ensureDir(this.promptsDir);
      
      const availableChatmates = await this.getAvailableChatmates();
      
      for (const agentName of agentNames) {
        const matchingFiles = availableChatmates.filter(file => 
          file.toLowerCase().includes(agentName.toLowerCase()) ||
          file === `${agentName}.chatmode.md`
        );

        if (matchingFiles.length === 0) {
          console.warn(chalk.yellow(`⚠️  No chatmate found matching: ${agentName}`));
          continue;
        }

        if (matchingFiles.length > 1) {
          console.warn(chalk.yellow(`⚠️  Multiple chatmates found for "${agentName}":`));
          matchingFiles.forEach(file => console.log(`    - ${file}`));
          console.log(chalk.yellow(`    Installing all matches...\n`));
        }

        for (const file of matchingFiles) {
          await this.installChatmate(file, force);
        }
      }

    } catch (error) {
      throw new Error(`Failed to install specific chatmates: ${error.message}`);
    }
  }

  /**
   * Install a single chatmate file
   * @param {string} filename - Chatmate filename to install
   * @param {boolean} force - Force reinstall even if already installed
   */
  async installChatmate(filename, force = false) {
    try {
      const sourcePath = path.join(this.matesDir, filename);
      const destPath = path.join(this.promptsDir, filename);

      // Check if source file exists
      if (!await fs.pathExists(sourcePath)) {
        throw new Error(`Chatmate file not found: ${filename}`);
      }

      // Check if already installed and not forcing
      if (!force && await fs.pathExists(destPath)) {
        console.log(chalk.gray(`⏭️  ${filename} (already installed)`));
        return;
      }

      // Copy file
      await fs.copy(sourcePath, destPath, { overwrite: true });
      
      const status = force ? 'reinstalled' : 'installed';
      console.log(chalk.green(`✅ ${filename} (${status})`));

    } catch (error) {
      console.error(chalk.red(`❌ ${filename} (failed: ${error.message})`));
      throw error;
    }
  }

  /**
   * Uninstall all chatmate agents
   */
  async uninstallAll() {
    try {
      const installedChatmates = await this.getInstalledChatmates();
      
      if (installedChatmates.length === 0) {
        console.log(chalk.yellow('No chatmates currently installed.'));
        return;
      }

      for (const chatmate of installedChatmates) {
        await this.uninstallChatmate(chatmate);
      }

    } catch (error) {
      throw new Error(`Failed to uninstall all chatmates: ${error.message}`);
    }
  }

  /**
   * Uninstall specific chatmate agents
   * @param {string[]} agentNames - Array of agent names to uninstall
   */
  async uninstallSpecific(agentNames) {
    try {
      const installedChatmates = await this.getInstalledChatmates();
      
      for (const agentName of agentNames) {
        const matchingFiles = installedChatmates.filter(file => 
          file.toLowerCase().includes(agentName.toLowerCase()) ||
          file === `${agentName}.chatmode.md`
        );

        if (matchingFiles.length === 0) {
          console.warn(chalk.yellow(`⚠️  No installed chatmate found matching: ${agentName}`));
          continue;
        }

        for (const file of matchingFiles) {
          await this.uninstallChatmate(file);
        }
      }

    } catch (error) {
      throw new Error(`Failed to uninstall specific chatmates: ${error.message}`);
    }
  }

  /**
   * Uninstall a single chatmate file
   * @param {string} filename - Chatmate filename to uninstall
   */
  async uninstallChatmate(filename) {
    try {
      const filePath = path.join(this.promptsDir, filename);

      if (!await fs.pathExists(filePath)) {
        console.warn(chalk.yellow(`⚠️  ${filename} (not installed)`));
        return;
      }

      await fs.remove(filePath);
      console.log(chalk.red(`🗑️  ${filename} (uninstalled)`));

    } catch (error) {
      console.error(chalk.red(`❌ ${filename} (failed to uninstall: ${error.message})`));
      throw error;
    }
  }

  /**
   * List available and installed chatmates
   * @param {Object} options - List options
   */
  async listChatmates(options = {}) {
    try {
      const availableChatmates = await this.getAvailableChatmates();
      const installedChatmates = await this.getInstalledChatmates();

      if (!options.installed) {
        console.log(chalk.blue('📦 Available Chatmates:'));
        if (availableChatmates.length === 0) {
          console.log(chalk.gray('  No chatmates available.'));
        } else {
          availableChatmates.forEach(chatmate => {
            const isInstalled = installedChatmates.includes(chatmate);
            const status = isInstalled ? chalk.green('✅ installed') : chalk.gray('⏸️  available');
            const name = chatmate.replace('.chatmode.md', '');
            console.log(`  ${name} ${status}`);
          });
        }
        console.log();
      }

      if (!options.available) {
        console.log(chalk.green('✅ Installed Chatmates:'));
        if (installedChatmates.length === 0) {
          console.log(chalk.gray('  No chatmates currently installed.'));
          console.log(chalk.gray('  Run "chatmate hire" to install all available chatmates.'));
        } else {
          installedChatmates.forEach(chatmate => {
            const name = chatmate.replace('.chatmode.md', '');
            console.log(`  ${name} ${chalk.green('✅ installed')}`);
          });
        }
        console.log();
      }

      // Summary
      console.log(chalk.cyan(`📊 Summary: ${installedChatmates.length}/${availableChatmates.length} chatmates installed`));

    } catch (error) {
      throw new Error(`Failed to list chatmates: ${error.message}`);
    }
  }

  /**
   * Show ChatMate and VS Code installation status
   */
  async showStatus() {
    try {
      console.log(chalk.blue('🔍 ChatMate Installation Status\n'));

      // Check VS Code installation
      const vsCodeInstalled = await this.checkVSCodeInstallation();
      const vsCodeStatus = vsCodeInstalled ? 
        chalk.green('✅ VS Code detected') : 
        chalk.red('❌ VS Code not found');
      console.log(`VS Code: ${vsCodeStatus}`);

      // Check prompts directory
      const promptsDirExists = await fs.pathExists(this.promptsDir);
      const promptsStatus = promptsDirExists ? 
        chalk.green('✅ Prompts directory exists') : 
        chalk.yellow('⚠️  Prompts directory not found');
      console.log(`Prompts Directory: ${promptsStatus}`);
      console.log(`Path: ${chalk.cyan(this.promptsDir)}\n`);

      // Show chatmate statistics
      const availableChatmates = await this.getAvailableChatmates();
      const installedChatmates = await this.getInstalledChatmates();
      
      console.log(chalk.blue('📊 Chatmate Statistics:'));
      console.log(`Available: ${chalk.cyan(availableChatmates.length)} chatmates`);
      console.log(`Installed: ${chalk.green(installedChatmates.length)} chatmates`);
      
      if (installedChatmates.length < availableChatmates.length) {
        const uninstalled = availableChatmates.length - installedChatmates.length;
        console.log(chalk.yellow(`Pending: ${uninstalled} chatmates not installed`));
        console.log(chalk.gray('\nRun "chatmate hire" to install all available chatmates.'));
      }

    } catch (error) {
      throw new Error(`Failed to show status: ${error.message}`);
    }
  }

  /**
   * Check if VS Code is installed
   * @returns {Promise<boolean>} True if VS Code is installed
   */
  async checkVSCodeInstallation() {
    const platform = os.platform();
    
    try {
      switch (platform) {
        case 'darwin': // macOS
          return await fs.pathExists('/Applications/Visual Studio Code.app');
        
        case 'linux': // Linux
          // Check common installation paths
          const linuxPaths = [
            '/usr/bin/code',
            '/usr/local/bin/code',
            '/snap/bin/code'
          ];
          for (const codePath of linuxPaths) {
            if (await fs.pathExists(codePath)) return true;
          }
          return false;
        
        case 'win32': // Windows
          const programFiles = process.env.PROGRAMFILES || 'C:\\Program Files';
          const programFilesX86 = process.env['PROGRAMFILES(X86)'] || 'C:\\Program Files (x86)';
          const winPaths = [
            path.join(programFiles, 'Microsoft VS Code', 'Code.exe'),
            path.join(programFilesX86, 'Microsoft VS Code', 'Code.exe')
          ];
          for (const codePath of winPaths) {
            if (await fs.pathExists(codePath)) return true;
          }
          return false;
        
        default:
          return false;
      }
    } catch (error) {
      return false;
    }
  }

  /**
   * Manage ChatMate configuration settings
   * @param {Object} options - Configuration options
   */
  async manageConfig(options = {}) {
    try {
      if (options.show) {
        console.log(chalk.blue('⚙️  ChatMate Configuration:\n'));
        console.log(`Mates Directory: ${chalk.cyan(this.matesDir)}`);
        console.log(`Prompts Directory: ${chalk.cyan(this.promptsDir)}`);
        console.log(`Platform: ${chalk.cyan(os.platform())}`);
        console.log(`Node Version: ${chalk.cyan(process.version)}`);
        return;
      }

      if (options.reset) {
        console.log(chalk.yellow('🔄 Configuration reset is not implemented yet.'));
        console.log(chalk.gray('Current configuration is derived from system paths.'));
        return;
      }

      console.log(chalk.blue('⚙️  Configuration Management:'));
      console.log('  --show    Show current configuration');
      console.log('  --reset   Reset configuration to defaults');

    } catch (error) {
      throw new Error(`Failed to manage configuration: ${error.message}`);
    }
  }
}

module.exports = { ChatMateManager };
