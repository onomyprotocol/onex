use onomy_test_lib::dockerfiles::{onomy_std_cosmos_daemon_with_arbitrary, ONOMYD_VERSION};

pub const PROVIDER_VERSION: &str = ONOMYD_VERSION;
pub const CONSUMER_ID: &str = "onex";
pub const CONSUMER_VERSION: &str = "v0.1.0";
pub const PROVIDER_ACCOUNT_PREFIX: &str = "onomy";
pub const CONSUMER_ACCOUNT_PREFIX: &str = "onomy";

pub fn consumer_binary_name() -> String {
    format!("{CONSUMER_ID}d")
}

pub fn consumer_directory() -> String {
    //format!(".{CONSUMER_ID}")
    format!(".onomy_{CONSUMER_ID}")
}

#[rustfmt::skip]
const DOWNLOAD_ONOMYD: &str = r#"ADD https://github.com/onomyprotocol/onomy/releases/download/$DAEMON_VERSION/onomyd $DAEMON_HOME/cosmovisor/genesis/$DAEMON_VERSION/bin/onomyd"#;

pub fn dockerfile_onomyd() -> String {
    onomy_std_cosmos_daemon_with_arbitrary("onomyd", ".onomy", ONOMYD_VERSION, DOWNLOAD_ONOMYD)
}
