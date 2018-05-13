## simple sentence counter counter
read stdin and count it

### integrated with FZF to fast walk through cd history
inspired by http://blog.naichilab.com/entry/zsh-percol

```
typeset -U chpwd_functions
## log cd history
CD_HISTORY_FILE=${HOME}/.cd_history_file # cd history filcd history file
function chpwd_record_history() {
    echo $PWD | sentence-counter -dest=${CD_HISTORY_FILE}
}
chpwd_functions=($chpwd_functions chpwd_record_history)


function fd(){
    dest=$(sentence-counter -dest=${CD_HISTORY_FILE} --show-reverse |  fzf +m --query "$LBUFFER" --prompt="cd > ")
		cd $dest
}


```
