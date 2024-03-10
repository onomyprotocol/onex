// for locally cleaning up log files and other temporaries

use onomy_test_lib::super_orchestrator::{remove_files_in_dir, stacked_errors::Result, std_init};

#[tokio::main]
async fn main() -> Result<()> {
    std_init()?;

    // Similar to the top level .gitignore, we catch all things that will be used on
    // all branches
    remove_files_in_dir("./", &["consumer-democracy"]).await?;
    remove_files_in_dir("./tests/dockerfiles", &["__tmp.dockerfile"]).await?;
    remove_files_in_dir("./tests/dockerfiles/dockerfile_resources", &[
        "onomyd",
        "onexd",
        "appnamed",
        "__tmp_hermes_config.toml",
        "havend",
        "arc_ethd",
    ])
    .await?;
    remove_files_in_dir("./tests/logs", &[".log", ".json", ".toml"]).await?;
    remove_files_in_dir("./tests/resources/keyring-test/", &[".address", ".info"]).await?;

    Ok(())
}
