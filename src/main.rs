use std::env;
use std::process;
use std::io;

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
                130 // because when I press ctrl-c and echo $? right after, it prints 130
            },
        },
        Err(e) => {
            println!("Error executing command: {}", e);
            2
        },
    });
}

fn print_usage() {
    println!("Usage: ffind [-fdri] [-e=maxdepth] [--debug --help] [expression] [path]");
}

fn make_command(args_in: Vec<String>) -> Result<Vec<String>, io::Error> {
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

fn sort_args(args_in: Vec<String>) -> Result<SortedArgs, io::Error> {
    let mut short_args = Vec::new();
    let mut long_args = Vec::new();
    let mut file_name = String::from("");
    let mut path = String::from("");
    let mut exec_args = String::from("");

    short_args = get_short_args(args_in)?;
    long_args = get_long_args(args_in)?;

    let mut sorted_args = SortedArgs {
        short_args: short_args,
        long_args: long_args,
        file_name: file_name,
        path: path,
        exec_args: exec_args,
    };

    // Err(io::Error::new(io::ErrorKind::Other, "testing"))
    Ok(sorted_args)
}

fn get_short_args(args_in: Vec<String>) -> Result<Vec<String>, io::Error> {
    let mut ret_short_args: Vec<String> = Vec::new();

    for args in args_in {
        let args_c = args.char_indices();
        if args_c.count() > 1 { // short flags or long flags
            if args_c.next().unwrap() == (0, '-') {
                let second_args_char = args_c.next().unwrap();
                if second_args_char == (1, '-') { // second dash so long flag
                    continue;
                } else {
                    ret_short_args.push(second_args_char.1.to_string());
                    return Ok(ret_short_args);
                }
            }
        }
    }
    return Ok(ret_short_args);
}

fn get_long_args(args_in: Vec<String>) -> Result<Vec<String>, io::Error> {
    let mut ret_long_args: Vec<String> = Vec::new();
    Ok(ret_long_args)
}
