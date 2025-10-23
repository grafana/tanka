mod cli;
mod environment;
mod export;
mod jsonnet;
mod manifest;
mod native;
mod template;

use anyhow::Result;
use clap::Parser;
use cli::Cli;
use env_logger::Env;

#[tokio::main]
async fn main() -> Result<()> {
    env_logger::Builder::from_env(Env::default().default_filter_or("info")).init();

    let cli = Cli::parse();
    cli.run().await
}
