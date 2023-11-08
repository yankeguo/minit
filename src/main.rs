mod env;
mod logger;
mod unit;

use crate::env::{env_bool, env_dir};
use crate::unit::load_units;
use std::error::Error;
use std::fs;

fn main() -> Result<(), Box<dyn Error>> {
    let opt_dir_unit = env_dir("MINIT_UNIT_DIR", Some(String::from("/etc/minit.d")));
    let opt_dir_log = env_dir("MINIT_LOG_DIR", Some(String::from("/var/log/minit")));
    let opt_quick_exit = env_bool("MINIT_QUICK_EXIT");

    match &opt_dir_unit {
        None => {}
        Some(dir) => fs::create_dir_all(dir).expect("create minit directory"),
    }

    match &opt_dir_log {
        None => {}
        Some(dir) => fs::create_dir_all(dir).expect("create log directory"),
    }

    println!("minit: starting (#{})", env!("MINIT_COMMIT"));

    let units = match &opt_dir_unit {
        None => vec![],
        Some(dir) => load_units(dir)?,
    };

    if units.is_empty() && opt_quick_exit {
        return Ok(());
    }

    println!("units loaded {:?}", units.len());

    Ok(())
}
