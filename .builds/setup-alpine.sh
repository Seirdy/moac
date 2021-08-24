#!/bin/sh
set -e

rm -f /usr/bin/ld /usr/bin/cc
ln -sf /usr/bin/ld.lld /usr/bin/ld
ln -sf /usr/bin/clang-1* /usr/bin/cc
cd /usr/lib/clang/*
mkdir -p share
echo "[cfi-unrelated-cast]
# The specification of std::get_temporary_buffer mandates a cast to
# uninitialized T* (libstdc++, MSVC stdlib).
fun:_ZSt20get_temporary_buffer*
fun:*get_temporary_buffer@.*@std@@*

# STL address-of magic (libstdc++).
fun:*__addressof*

# Windows C++ stdlib headers that contain bad unrelated casts.
src:*xmemory0
src:*xstddef

# std::_Sp_counted_ptr_inplace::_Sp_counted_ptr_inplace() (libstdc++).
# This ctor is used by std::make_shared and needs to cast to uninitialized T*
# in order to call std::allocator_traits<T>::construct.
fun:_ZNSt23_Sp_counted_ptr_inplace*" >share/cfi_blacklist.txt
