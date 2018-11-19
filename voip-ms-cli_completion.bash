#/usr/bin/env bash
# https://iridakos.com/tutorials/2018/03/01/bash-programmable-completion-tutorial.html
_voipmscli_completions()
{
  COMPREPLY=($(compgen -W "show-balance show-recent block-recent block-number" "${COMP_WORDS[1]}"))
}

complete -F _voipmscli_completions voip-ms-cli
