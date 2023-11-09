use std::error::Error;

pub type GeneralResult<T> = Result<T, Box<dyn Error>>;
