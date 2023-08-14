mod env;
mod logger;

use crate::env::{env_bool, env_str};
use std::fs;

fn main() {
    let mut opt_dir_unit: Option<String> = Some(String::from("/etc/minit.d"));
    let mut opt_dir_log: Option<String> = None;
    let mut opt_quick_exit = false;

    env_str("MINIT_UNIT_DIR", &mut opt_dir_unit);
    env_str("MINIT_LOG_DIR", &mut opt_dir_log);
    env_bool("MINIT_QUICK_EXIT", &mut opt_quick_exit);

    if opt_dir_log.is_some() {
        let opt_dir_log = opt_dir_log.unwrap();
        fs::create_dir_all(opt_dir_log).expect("create log directory");
    }

    if opt_dir_unit.is_some() {
        let _opt_dir_unit = opt_dir_unit.unwrap();
    }

    println!("minit: starting (#{})", env!("MINIT_COMMIT"));
}
