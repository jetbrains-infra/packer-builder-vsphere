<#
.SYNOPSIS
  A flixible build script for PowerShell
.DESCRIPTION
  This build script takes parameters to make it easier to only build one of the products if you are only interested in one, but will default to buildign all.
.EXAMPLE
  & build.ps1
  This command will build both tools for all three operating systems
.EXAMPLE
  & .\build.ps1 -builder iso -GOOS windows -verbose
  This command will build only the ISO builder for windows, and print the verbose output from go build
.EXAMPLE
  & .\build.ps1 -builder iso -GOOS windows,linux -verbose
  This command builds the iso builder for windows and linux.
.NOTES
  This script should be fully cross platform compatible with Windows, Linux and Mac.
#>
[cmdletbinding()]
param(
  # Would you like to build the ISO plugin or the Clone Plugin or both?
  [ValidateSet('iso','clone','all')]
  [string[]]$builder = 'all',
  # Value for the CGO_ENABLED Env Var. Default to 0
  [int]$CGO_ENABLED = 0,
  # Value for the GOARCH Env Var. Default to amd64
  [string]$GOARCH = 'amd64',
  # Value for the GOOS env Var. Which OS would you like to build for? Default to All.
  [ValidateSet('darwin','linux','windows','all')]
  [string[]]$GOOS = 'all',
  # Destination Directory for artifacts
  [string]$dest = './bin'
)

function Set-EnvVarValue {
  param(
    $varname,
    $value
  )

  if(Test-Path "env:\$varname"){
    Set-Item "env:\$varname" -Value $value
  } else {
    New-Item "env:\$varname" -Value $value
  }
}

if(-not (Test-Path $dest)) {
  New-Item -Path $dest -ItemType Directory
}

Remove-Item $dest\* -Recurse -Force

$env_vars = @('CGO_ENABLED','GOARCH','GOOS')

foreach ($env_var in $env_vars) {
  Set-EnvVarValue -varname $env_var -value (Get-Variable $env_var).value
}

if($GOOS -eq 'all') {
  $OSList = @('darwin','linux','windows')
} else {
  $OSList = $GOOS
}

if($builder -eq 'all') {
  $builderList = @('iso','clone')
} else {
  $builderList = $builder
}

foreach($os in $OSList){
  foreach($type in $builderList){
    Set-EnvVarValue -varname 'GOOS' -value $os
    $command = 'go build '

    if($VerbosePreference -eq 'continue') {
      $command += '-v '
    }

    $command += "-o $dest/packer-builder-vsphere-$type.$os ./cmd/$type"
    write-Verbose $command
    Invoke-Expression -Command $command
  }
}

foreach ($env_var in $env_vars) {
  Remove-Item "env:\$env_var"
}
