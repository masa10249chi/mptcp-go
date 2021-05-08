set encoding=utf-8
scriptencoding utf-8

set mouse-=a
set number
set autoindent
set autowrite
set tabstop=4
set shiftwidth=4
set expandtab

call plug#begin()
Plug 'fatih/vim-go'
Plug 'fatih/molokai'
call plug#end()

"Plug 'fatih/vim-go'
let g:go_version_warning = 0
let g:go_fmt_autosave = 0

"Plug 'fatih/molokai'
let g:go_highlight_functions = 1
let g:go_highlight_methods = 1
let g:go_highlight_structs = 1
let g:rehash256 = 1
let g:molokai_original = 1
colorscheme molokai
syntax on

"Plug 'nathanaelkane/vim-indent-guides'
let g:indent_guides_enable_on_vim_startup = 1

autocmd FileType go nmap <leader>b  <Plug>(go-build)
autocmd FileType go nmap <leader>r  <Plug>(go-run)