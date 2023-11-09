mod logger;
mod types;
mod unit;
mod utils;

use crate::types::GeneralResult;
use crate::unit::{load_arg_unit, load_dir_units, load_env_unit};
use std::str::FromStr;
use std::{env, fs};

pub fn env_dir(key: &str, default_value: Option<String>) -> Option<String> {
    match env::var(key) {
        Ok(val) if val.eq_ignore_ascii_case("none") => None,
        Ok(val) if val.is_empty() => default_value,
        Ok(val) => Some(val),
        _ => None,
    }
}

pub fn env_bool(key: &str) -> bool {
    match env::var(key) {
        Ok(val) => match <bool as FromStr>::from_str(val.as_str()) {
            Ok(val) => val,
            _ => false,
        },
        _ => false,
    }
}

fn main() -> GeneralResult<()> {
    let opt_dir_unit = env_dir("MINIT_UNIT_DIR", Some(String::from("/etc/minit.d")));
    let opt_dir_log = env_dir("MINIT_LOG_DIR", Some(String::from("/var/log/minit")));
    let opt_quick_exit = env_bool("MINIT_QUICK_EXIT");

    if let Some(dir) = &opt_dir_log {
        fs::create_dir_all(dir)?;
    }
    if let Some(dir) = &opt_dir_unit {
        fs::create_dir_all(dir)?;
    }

    println!("minit: starting (#{})", env!("MINIT_COMMIT"));

    let mut units = vec![];

    if let Some(dir) = &opt_dir_unit {
        load_dir_units(&mut units, dir)?;
    }
    load_env_unit(&mut units, &mut env::vars())?;
    load_arg_unit(&mut units, &mut env::args())?;

    if units.is_empty() && opt_quick_exit {
        return Ok(());
    }

    println!("units loaded {:?}", units.len());

    Ok(())
}
