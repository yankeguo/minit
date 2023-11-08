use std::env;
use std::str::FromStr;

pub fn env_dir(key: &str, default_value: Option<String>) -> Option<String> {
    return match env::var(key) {
        Ok(val) => {
            if val.eq_ignore_ascii_case("none") {
                None
            } else if val.is_empty() {
                default_value
            } else {
                Some(val)
            }
        }
        _ => None,
    };
}

pub fn env_bool(key: &str) -> bool {
    return match env::var(key) {
        Ok(val) => match <bool as FromStr>::from_str(val.as_str()) {
            Ok(val) => val,
            _ => false,
        },
        _ => false,
    };
}
