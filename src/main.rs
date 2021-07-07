use std::env;
use std::process;
use std::error;
use std::io;
use std::fmt;

fn main() {
    let args: Vec<String> = env::args().collect();

    // making the command
    let cmd_args_result = make_command(args);
    let cmd_args = match cmd_args_result {
        Err(e) => {
            println!("Error making command: {}", e);
            process::exit(1);
        },
        Ok(cmd_args) => cmd_args,
    };

    // executing the command
    let exec_status_result = process::Command::new("find")
        .args(cmd_args)
        .stdout(process::Stdio::inherit())
        .stderr(process::Stdio::inherit())
        .status();

    // proper error code propogation
    // spaghettified because Command::status() returns - Ok(Exitstatus(Exitstatus(code)))
    // we need to unwrap highest Ok then unwrap the Exitstatus to get to the return code of
    // command executed
    process::exit(match exec_status_result {
        Ok(exec_out) => match exec_out.code() {
            Some(code) => code,
            None       => {
                println!("Process terminated by signal");
                130
            },
        },
        Err(e) => {
            println!("Error executing command: {}", e);
            1
        },
    });
}

fn print_usage() {
    println!("Usage: ffind [-fdri] [-e=maxdepth] [--debug --help] [expression] [path]");
}

#[derive(Debug)]
enum CommandMakingError {
    Error1,
    Error2,
}

impl fmt::Display for CommandMakingError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match *self {
            CommandMakingError::Error1 => write!(f,"1st command making error"),
            CommandMakingError::Error2 => write!(f,"2st command making error"),
        }
    }
}

impl error::Error for CommandMakingError {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        match *self {
            CommandMakingError::Error1    => None,
            CommandMakingError::Error2(ref e) => Some(e),
        }
    }
}

// to enable `?` called for CommandMakingError in a function returning ArgsSortingError
impl From<CommandMakingError> for ArgsSortingError {
    fn from(err: CommandMakingError) -> ArgsSortingError {
        ArgsSortingError::ArgsSortingError2(err)
    }
}

fn make_command(args_in: Vec<String>) -> Result<Vec<String>, CommandMakingError> {
    let mut args_out = Vec::new();

    let sorted_args = sort_args(args_in)?;
    Ok(args_out)
}

#[derive(Debug)]
struct SortedArgs {
    short_args: Vec<String>,
    long_args: Vec<String>,
    file_name: String,
    path: String,
    exec_args: String,
}

#[derive(Debug)]
enum ArgsSortingError {
    ArgsSortingError1,
    ArgsSortingError2,
}

impl fmt::Display for ArgsSortingError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match *self {
            ArgsSortingError::ArgsSortingError1 => write!(f,"1st args sorting"),
            ArgsSortingError::ArgsSortingError2 => write!(f,"2st args sorting"),
        }
    }
}

impl error::Error for ArgsSortingError {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        match *self {
            ArgsSortingError::ArgsSortingError1        => None,
            ArgsSortingError::ArgsSortingError2(ref e) => Some(e),
        }
    }
}

// to enable `?` called for ArgsSortingError in a function returning CommandSortingError
impl From<ArgsSortingError> for CommandMakingError {
    fn from(err: ArgsSortingError) -> CommandMakingError {
        CommandMakingError::ArgsSortingError2(err)
    }
}

fn sort_args(args_in: Vec<String>) -> Result<SortedArgs, ArgsSortingError> {
    let mut short_args = Vec::new();
    let mut long_args = Vec::new();
    let mut file_name = String::from("");
    let mut path = String::from("");
    let mut exec_args = String::from("");

    let mut sorted_args = SortedArgs {
        short_args: short_args,
        long_args: long_args,
        file_name: file_name,
        path: path,
        exec_args: exec_args,
    };

    // let long_args = longArgs(args_in)?;

    Err(ArgsSortingError::ArgsSortingError2)
    // Ok(sorted_args)
}
